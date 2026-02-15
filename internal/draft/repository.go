package draft

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type DraftRepository struct {
	rdb *redis.Client
}

func NewDraftRepository(rdb *redis.Client) *DraftRepository {
	return &DraftRepository{rdb: rdb}
}

func key(id string) string {
	return fmt.Sprintf("txn:draft:%s", id)
}

// ตอน save ให้ผูก user → draftID
func (r *DraftRepository) Save(ctx context.Context, d DraftTxn) error {

	k := key(d.DraftID)

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	pipe := r.rdb.TxPipeline()

	pipe.Set(ctx, k, b, 24*time.Hour)

	// 👇 index user → draftID
	userKey := fmt.Sprintf("txn:user:%s:drafts", d.UserID)
	pipe.SAdd(ctx, userKey, d.DraftID)
	pipe.Expire(ctx, userKey, 24*time.Hour)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *DraftRepository) ListByUser(ctx context.Context, userID string) ([]DraftTxn, error) {

	userKey := fmt.Sprintf("txn:user:%s:drafts", userID)

	ids, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	var result []DraftTxn

	for _, id := range ids {
		d, err := r.Get(ctx, id)

		if err == nil {
			result = append(result, *d)
		} else {
			// cleanup index ที่หมดอายุ
			r.rdb.SRem(ctx, userKey, id)
		}
	}

	return result, nil
}

func (r *DraftRepository) Get(ctx context.Context, draftID string) (*DraftTxn, error) {

	val, err := r.rdb.Get(ctx, key(draftID)).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("draft not found")
	}

	if err != nil {
		return nil, err
	}

	var d DraftTxn
	if err := json.Unmarshal([]byte(val), &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *DraftRepository) Delete(ctx context.Context, draftID string) error {

	d, err := r.Get(ctx, draftID)

	pipe := r.rdb.TxPipeline()

	if err == nil {
		userKey := fmt.Sprintf("txn:user:%s:drafts", d.UserID)
		pipe.SRem(ctx, userKey, draftID)
	}

	pipe.Del(ctx, key(draftID))

	_, err = pipe.Exec(ctx)
	return err
}

func (r *DraftRepository) UpdateStatus(ctx context.Context, draftID string, status DraftStatus, errMsg string) error {
	d, err := r.Get(ctx, draftID)
	if err != nil {
		return err
	}

	d.Status = status
	d.Error = errMsg
	d.UpdatedAt = time.Now()

	return r.Save(ctx, *d)
}

func (r *DraftRepository) StatsByUser(ctx context.Context, userID string) (*DraftStats, error) {

	userKey := fmt.Sprintf("txn:user:%s:drafts", userID)

	ids, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	stats := &DraftStats{}

	for _, id := range ids {

		d, err := r.Get(ctx, id)
		if err != nil {
			// cleanup orphan index
			r.rdb.SRem(ctx, userKey, id)
			continue
		}

		stats.Total++

		switch d.Status {
		case DraftStatusQueued:
			stats.Queued++
		case DraftStatusProcessing:
			stats.Processing++
		case DraftStatusWaitingConfirm:
			stats.WaitingConfirm++
		case DraftStatusFailed:
			stats.Failed++
		}
	}

	return stats, nil
}

package draftrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Save(ctx context.Context, d model.DraftTxn) error
	ListByUser(ctx context.Context, userID string) ([]model.DraftTxn, error)
	Get(ctx context.Context, draftID string) (*model.DraftTxn, error)
	Delete(ctx context.Context, draftID string) error
	UpdateStatus(ctx context.Context, draftID string, status model.DraftStatus, errMsg string) error
	StatsByUser(ctx context.Context, userID string) (*model.DraftStats, error)
}
type repository struct {
	rdb *redis.Client
}

func NewDraftRepository(rdb *redis.Client) Repository {
	return &repository{rdb}
}

func key(id string) string {
	return fmt.Sprintf("txn:draft:%s", id)
}

// ตอน save ให้ผูก user → draftID
func (r *repository) Save(ctx context.Context, d model.DraftTxn) error {

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

func (r *repository) ListByUser(ctx context.Context, userID string) ([]model.DraftTxn, error) {

	userKey := fmt.Sprintf("txn:user:%s:drafts", userID)

	ids, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	var result []model.DraftTxn

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

func (r *repository) Get(ctx context.Context, draftID string) (*model.DraftTxn, error) {

	val, err := r.rdb.Get(ctx, key(draftID)).Result()

	if err == redis.Nil {
		return nil, fmt.Errorf("draft not found")
	}

	if err != nil {
		return nil, err
	}

	var d model.DraftTxn
	if err := json.Unmarshal([]byte(val), &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *repository) Delete(ctx context.Context, draftID string) error {

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

func (r *repository) UpdateStatus(ctx context.Context, draftID string, status model.DraftStatus, errMsg string) error {
	d, err := r.Get(ctx, draftID)
	if err != nil {
		return err
	}

	d.Status = status
	d.Error = errMsg
	d.UpdatedAt = time.Now()

	return r.Save(ctx, *d)
}

func (r *repository) StatsByUser(ctx context.Context, userID string) (*model.DraftStats, error) {

	userKey := fmt.Sprintf("txn:user:%s:drafts", userID)

	ids, err := r.rdb.SMembers(ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	stats := &model.DraftStats{}

	for _, id := range ids {

		d, err := r.Get(ctx, id)
		if err != nil {
			// cleanup orphan index
			r.rdb.SRem(ctx, userKey, id)
			continue
		}

		stats.Total++

		switch d.Status {
		case model.DraftStatusQueued:
			stats.Queued++
		case model.DraftStatusProcessing:
			stats.Processing++
		case model.DraftStatusWaitingConfirm:
			stats.WaitingConfirm++
		case model.DraftStatusFailed:
			stats.Failed++
		}
	}

	return stats, nil
}

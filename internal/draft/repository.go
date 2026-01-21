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

// ตอน save ให้ผูก user → traceID
func (r *DraftRepository) Save(ctx context.Context, d DraftTxn) error {

	k := key(d.TraceID)

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	pipe := r.rdb.TxPipeline()

	pipe.Set(ctx, k, b, 24*time.Hour)

	// 👇 index user → traceID
	userKey := fmt.Sprintf("txn:user:%s:drafts", d.UserID)
	pipe.SAdd(ctx, userKey, d.TraceID)
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

func (r *DraftRepository) Get(ctx context.Context, traceID string) (*DraftTxn, error) {

	val, err := r.rdb.Get(ctx, key(traceID)).Result()

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

func (r *DraftRepository) Delete(ctx context.Context, traceID string) error {

	d, err := r.Get(ctx, traceID)

	pipe := r.rdb.TxPipeline()

	if err == nil {
		userKey := fmt.Sprintf("txn:user:%s:drafts", d.UserID)
		pipe.SRem(ctx, userKey, traceID)
	}

	pipe.Del(ctx, key(traceID))

	_, err = pipe.Exec(ctx)
	return err
}


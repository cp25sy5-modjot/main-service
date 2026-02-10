package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const TaskPurgeUser = "user:purge"

type PurgeUserPayload struct {
	UserID string `json:"user_id"`
}

func NewPurgeUserTask(userID string, delaySeconds int64) (*asynq.Task, error) {
	payload, err := json.Marshal(PurgeUserPayload{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(
		TaskPurgeUser,
		payload,
		asynq.ProcessIn(time.Duration(delaySeconds)*time.Second),
	), nil
}

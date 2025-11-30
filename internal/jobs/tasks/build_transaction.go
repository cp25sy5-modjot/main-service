package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TaskBuildTransaction = "transaction:build_from_image"

type BuildTransactionPayload struct {
	UserID  string `json:"user_id"`
	Path    string `json:"path"`
	TraceID string `json:"trace_id"`
}

func NewBuildTransactionTask(userID, path, traceID string) (*asynq.Task, error) {
	payload, err := json.Marshal(BuildTransactionPayload{
		UserID:  userID,
		Path:    path,
		TraceID: traceID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskBuildTransaction, payload), nil
}

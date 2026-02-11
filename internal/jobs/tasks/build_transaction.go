package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TaskBuildTransaction = "transaction:build_from_image"

type BuildTransactionPayload struct {
	UserID  string `json:"user_id"`
	Path    string `json:"path"`
	DraftID string `json:"draft_id"`
}

func NewBuildTransactionTask(userID, path, draftID string) (*asynq.Task, error) {
	payload, err := json.Marshal(BuildTransactionPayload{
		UserID:  userID,
		Path:    path,
		DraftID: draftID,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TaskBuildTransaction, payload), nil
}

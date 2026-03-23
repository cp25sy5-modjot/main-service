package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const TaskRunFixCost = "fixcost:run"

type RunFixCostPayload struct {
	FixCostID string `json:"fix_cost_id"`
	UserID    string `json:"user_id"`
}

func NewRunFixCostTask(fixCostID string, userID string) (*asynq.Task, error) {
	payload, err := json.Marshal(RunFixCostPayload{
		FixCostID: fixCostID,
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskRunFixCost, payload), nil
}

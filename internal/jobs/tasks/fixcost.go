package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const TypeProcessFixCost = "fixcost:process"

type ProcessFixCostPayload struct {
	FixCostID string    `json:"fix_cost_id"`
	Date      time.Time `json:"date"`
	UserID    string    `json:"user_id"`
}

func NewProcessFixCostTask(fcID string, date time.Time, userID string) (*asynq.Task, error) {
	payload, err := json.Marshal(ProcessFixCostPayload{
		FixCostID: fcID,
		Date:      date.UTC(),
		UserID:    userID,
	})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeProcessFixCost, payload), nil
}

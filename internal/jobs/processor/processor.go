package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hibiken/asynq"

	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
)

type Processor struct {
	txService txsvc.Service
	storage   storage.Storage
}

func NewProcessor(txService txsvc.Service, st storage.Storage) *Processor {
	return &Processor{
		txService: txService,
		storage:   st,
	}
}

func (p *Processor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(tasks.TaskBuildTransaction, p.handleBuildTransactionTask)
}

func (p *Processor) handleBuildTransactionTask(ctx context.Context, t *asynq.Task) error {
	var payload tasks.BuildTransactionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("JOB decode payload error: %v", err)
		return err
	}

	log.Printf("[JOB %s] Start transaction build. user=%s path=%s", payload.TraceID, payload.UserID, payload.Path)

	// 1. Load file
	data, err := p.storage.Load(ctx, payload.Path)
	if err != nil {
		log.Printf("[JOB %s] load error: %v", payload.TraceID, err)
		return err
	}

	// 2. Use your existing business logic
	tx, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {
		log.Printf("[JOB %s] process error: %v", payload.TraceID, err)
		return err
	}

	// 3. Optional: delete file
	if err := p.storage.Delete(ctx, payload.Path); err != nil {
		log.Printf("[JOB %s] delete file error: %v", payload.TraceID, err)
	}

	// 4. TODO: push notification / webhook here
	_ = tx

	log.Printf("[JOB %s] Done. transaction_id=%s", payload.TraceID, tx.TransactionID)
	return nil
}

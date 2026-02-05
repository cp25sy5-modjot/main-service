package processor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"

	d "github.com/cp25sy5-modjot/main-service/internal/draft"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
)

type Processor struct {
	txService txsvc.Service
	storage   storage.Storage

	draftRepo *d.DraftRepository
}

func NewProcessor(
	txService txsvc.Service,
	st storage.Storage,
	dr *d.DraftRepository,
) *Processor {
	return &Processor{
		txService: txService,
		storage:   st,
		draftRepo: dr,
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

	log.Printf("[JOB %s] Start transaction build. user=%s path=%s",
		payload.TraceID, payload.UserID, payload.Path)

	// 1. Load file
	data, err := p.storage.Load(ctx, payload.Path)
	if err != nil {
		log.Printf("[JOB %s] load error: %v", payload.TraceID, err)
		return err
	}

	// ===== STEP 1: mark processing =====
	_ = p.draftRepo.Save(ctx, d.DraftTxn{
		TraceID:   payload.TraceID,
		UserID:    payload.UserID,
		Status:    d.DraftStatusProcessing,
		UpdatedAt: time.Now(),
	})

	// ===== STEP 2: call AI =====
	draft, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {

		_ = p.draftRepo.Save(ctx, d.DraftTxn{
			TraceID:   payload.TraceID,
			UserID:    payload.UserID,
			Status:    d.DraftStatusFailed,
			Error:     err.Error(),
			UpdatedAt: time.Now(),
		})

		return err
	}

	// ===== STEP 3: save result to redis =====
	_ = p.draftRepo.Save(ctx, d.DraftTxn{
		TraceID: payload.TraceID,
		UserID:  payload.UserID,

		Status: d.DraftStatusWaitingConfirm,

		Title:     draft.Title,
		Date:      draft.Date,
		Items:     draft.Items,
		UpdatedAt: time.Now(),
	})

	// 4. delete file
	if err := p.storage.Delete(ctx, payload.Path); err != nil {
		log.Printf("[JOB %s] delete file error: %v", payload.TraceID, err)
	}

	log.Printf("[JOB %s] Done → waiting user confirm", payload.TraceID)

	return nil
}

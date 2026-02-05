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
	_ = p.draftRepo.UpdateStatus(ctx, payload.TraceID, d.DraftStatusProcessing, "")

	// ===== STEP 2: call AI =====
	draft, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {

		_ = p.draftRepo.UpdateStatus(ctx, payload.TraceID, d.DraftStatusFailed, err.Error())

		return err
	}
	exDraft, err := p.draftRepo.Get(ctx, payload.TraceID)
	if err != nil {
		log.Printf("[JOB %s] get draft error: %v", payload.TraceID, err)
		return err
	}

	exDraft.Title = draft.Title
	exDraft.Date = draft.Date
	exDraft.Items = draft.Items

	exDraft.Status = d.DraftStatusWaitingConfirm
	exDraft.UpdatedAt = time.Now()

	// ===== STEP 3: save result to redis =====
	_ = p.draftRepo.Save(ctx, *exDraft)

	// 4. delete file
	if err := p.storage.Delete(ctx, payload.Path); err != nil {
		log.Printf("[JOB %s] delete file error: %v", payload.TraceID, err)
	}

	log.Printf("[JOB %s] Done → waiting user confirm", payload.TraceID)

	return nil
}

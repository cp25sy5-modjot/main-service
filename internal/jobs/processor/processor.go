package processor

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hibiken/asynq"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	d "github.com/cp25sy5-modjot/main-service/internal/draft"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
)

type Processor struct {
	txService txsvc.Service
	storage   storage.Storage

	draftRepo *d.DraftRepository
	userRepo  *userrepo.Repository
}

func NewProcessor(
	txService txsvc.Service,
	st storage.Storage,
	dr *d.DraftRepository,
	userRepo *userrepo.Repository,
) *Processor {
	return &Processor{
		txService: txService,
		storage:   st,
		draftRepo: dr,
		userRepo:  userRepo,
	}
}

func (p *Processor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(tasks.TaskBuildTransaction, p.handleBuildTransactionTask)
	mux.HandleFunc(tasks.TaskPurgeUser, p.HandlePurgeUser)

}

func (p *Processor) handleBuildTransactionTask(ctx context.Context, t *asynq.Task) error {

	var payload tasks.BuildTransactionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("JOB decode payload error: %v", err)
		return err
	}

	log.Printf("[JOB %s] Start transaction build. user=%s path=%s",
		payload.DraftID, payload.UserID, payload.Path)

	// 1. Load file
	data, err := p.storage.Load(ctx, payload.Path)
	if err != nil {
		log.Printf("[JOB %s] load error: %v", payload.DraftID, err)
		return err
	}

	// ===== STEP 1: mark processing =====
	_ = p.draftRepo.UpdateStatus(ctx, payload.DraftID, d.DraftStatusProcessing, "")

	// ===== STEP 2: call AI =====
	draft, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {

		_ = p.draftRepo.UpdateStatus(ctx, payload.DraftID, d.DraftStatusFailed, err.Error())

		return err
	}
	exDraft, err := p.draftRepo.Get(ctx, payload.DraftID)
	if err != nil {
		log.Printf("[JOB %s] get draft error: %v", payload.DraftID, err)
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
	// if err := p.storage.Delete(ctx, payload.Path); err != nil {
	// 	log.Printf("[JOB %s] delete file error: %v", payload.DraftID, err)
	// }

	log.Printf("[JOB %s] Done → waiting user confirm", payload.DraftID)

	return nil
}

func (p *Processor) HandlePurgeUser(ctx context.Context, t *asynq.Task) error {
	var payload struct {
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("[JOB purge_user] payload error: %v", err)
		return err
	}

	log.Printf("[JOB purge_user] start user=%s", payload.UserID)

	user, err := p.userRepo.FindByID(payload.UserID)
	if err != nil {
		log.Printf("[JOB purge_user] user not found, skip")
		return nil
	}

	if user.Status != e.StatusInactive || user.UnsubscribedAt == nil {
		log.Printf("[JOB purge_user] user restored, skip")
		return nil
	}

	err = p.userRepo.HardDelete(payload.UserID)
	if err != nil {
		log.Printf("[JOB purge_user] delete error: %v", err)
		return err
	}

	log.Printf("[JOB purge_user] purged user=%s", payload.UserID)
	return nil
}

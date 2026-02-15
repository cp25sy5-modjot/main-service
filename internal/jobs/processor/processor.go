package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
func (p *Processor) handleBuildTransactionTask(ctx context.Context, t *asynq.Task) (err error) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[JOB] PANIC: %v", r)
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	var payload tasks.BuildTransactionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("[JOB] decode payload error: %+v", err)
		return err
	}

	start := time.Now()

	log.Printf("[JOB %s] START user=%s path=%s",
		payload.DraftID, payload.UserID, payload.Path)

	defer func() {
		log.Printf("[JOB %s] END in %s err=%v",
			payload.DraftID, time.Since(start), err)
	}()

	// --------------------------------------------------
	// 1) GET DRAFT
	// --------------------------------------------------
	exDraft, err := p.draftRepo.Get(ctx, payload.DraftID)
	if err != nil {
		log.Printf("[JOB %s] get draft error: %+v", payload.DraftID, err)
		return err
	}

	// ✅ IDEMPOTENT GUARD
	if exDraft.Status == d.DraftStatusWaitingConfirm {
		log.Printf("[JOB %s] already done → skip", payload.DraftID)
		return nil
	}

	// --------------------------------------------------
	// 2) LOAD FILE
	// --------------------------------------------------
	data, err := p.storage.Load(ctx, payload.Path)
	if err != nil {

		if os.IsNotExist(err) {

			retryCount, _ := asynq.GetRetryCount(ctx)
			maxRetry, _ := asynq.GetMaxRetry(ctx)

			log.Printf("[JOB %s] file not ready (%d/%d)",
				payload.DraftID, retryCount+1, maxRetry)

			if retryCount+1 >= maxRetry {

				_ = p.draftRepo.UpdateStatus(
					ctx,
					payload.DraftID,
					d.DraftStatusFailed,
					"file missing",
				)

				return asynq.SkipRetry
			}

			return fmt.Errorf("file not ready: %w", err)
		}

		log.Printf("[JOB %s] load file error: %+v", payload.DraftID, err)
		return err
	}

	// --------------------------------------------------
	// 3) UPDATE → PROCESSING
	// --------------------------------------------------
	if exDraft.Status != d.DraftStatusProcessing {
		if err := p.draftRepo.UpdateStatus(
			ctx,
			payload.DraftID,
			d.DraftStatusProcessing,
			"",
		); err != nil {
			log.Printf("[JOB %s] update processing status error: %+v",
				payload.DraftID, err)
		}
	}

	// --------------------------------------------------
	// 4) CALL AI
	// --------------------------------------------------
	log.Printf("[JOB %s] calling AI...", payload.DraftID)

	draftResult, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {

		log.Printf("[JOB %s] AI error: %+v", payload.DraftID, err)

		_ = p.draftRepo.UpdateStatus(
			ctx,
			payload.DraftID,
			d.DraftStatusFailed,
			err.Error(),
		)

		return err
	}

	log.Printf("[JOB %s] AI SUCCESS title=%s",
		payload.DraftID, draftResult.Title)

	// --------------------------------------------------
	// 5) SAVE RESULT
	// --------------------------------------------------
	exDraft.Title = draftResult.Title
	exDraft.Date = draftResult.Date
	exDraft.Items = draftResult.Items
	exDraft.Status = d.DraftStatusWaitingConfirm
	exDraft.UpdatedAt = time.Now()

	if err := p.draftRepo.Save(ctx, *exDraft); err != nil {
		log.Printf("[JOB %s] save draft error: %+v", payload.DraftID, err)

		// ❗ AI สำเร็จแล้ว → ห้าม retry ยิง AI ซ้ำ
		return asynq.SkipRetry
	}

	log.Printf("[JOB %s] DONE → waiting confirm", payload.DraftID)

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

	if user.Status != e.UserStatusInactive || user.UnsubscribedAt == nil {
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

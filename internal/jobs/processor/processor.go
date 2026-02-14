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

func (p *Processor) handleBuildTransactionTask(ctx context.Context, t *asynq.Task) error {

	var payload tasks.BuildTransactionPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("JOB decode payload error: %v", err)
		return err
	}

	log.Printf("[JOB %s] Start transaction build. user=%s path=%s",
		payload.DraftID, payload.UserID, payload.Path)

	data, err := p.storage.Load(ctx, payload.Path)
	if err != nil {

		if os.IsNotExist(err) {

			retryCount, ok1 := asynq.GetRetryCount(ctx)
			maxRetry, ok2 := asynq.GetMaxRetry(ctx)

			if !ok1 || !ok2 {
				log.Printf("[JOB %s] retry metadata missing", payload.DraftID)
			}

			log.Printf("[JOB %s] file not ready (retry %d/%d)",
				payload.DraftID, retryCount+1, maxRetry)

			if retryCount+1 >= maxRetry {

				if err := p.draftRepo.UpdateStatus(
					ctx,
					payload.DraftID,
					d.DraftStatusFailed,
					"file missing",
				); err != nil {
					log.Printf("update failed status error: %v", err)
				}

				return asynq.SkipRetry
			}

			return fmt.Errorf("file not ready: %w", err)
		}

		return err
	}

	if err := p.draftRepo.UpdateStatus(
		ctx,
		payload.DraftID,
		d.DraftStatusProcessing,
		"",
	); err != nil {
		log.Printf("update status error: %v", err)
	}

	draftResult, err := p.txService.ProcessUploadedFile(data, payload.UserID)
	if err != nil {

		if err := p.draftRepo.UpdateStatus(
			ctx,
			payload.DraftID,
			d.DraftStatusFailed,
			err.Error(),
		); err != nil {
			log.Printf("update failed status error: %v", err)
		}

		return err
	}

	exDraft, err := p.draftRepo.Get(ctx, payload.DraftID)
	if err != nil {
		log.Printf("[JOB %s] get draft error: %v", payload.DraftID, err)
		return err
	}

	exDraft.Title = draftResult.Title
	exDraft.Date = draftResult.Date
	exDraft.Items = draftResult.Items
	exDraft.Status = d.DraftStatusWaitingConfirm
	exDraft.UpdatedAt = time.Now()

	if err := p.draftRepo.Save(ctx, *exDraft); err != nil {
		log.Printf("[JOB %s] save draft error: %v", payload.DraftID, err)
		return err
	}

	/*
		if err := p.storage.Delete(ctx, payload.Path); err != nil {
			log.Printf("[JOB %s] delete file error: %v", payload.DraftID, err)
		}
	*/

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

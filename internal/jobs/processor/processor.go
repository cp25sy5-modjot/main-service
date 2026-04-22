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
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	drepo "github.com/cp25sy5-modjot/main-service/internal/draft/repository"
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	fcsvc "github.com/cp25sy5-modjot/main-service/internal/fix_cost/service"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
)

type Processor struct {
	txService txsvc.Service
	storage   storage.Storage

	draftRepo   drepo.Repository
	userRepo    userrepo.Repository
	fixCostRepo fcrepo.Repository
	client      *asynq.Client
	txRepo      txrepo.Repository
}

func NewProcessor(
	txService txsvc.Service,
	st storage.Storage,
	dr drepo.Repository,
	userRepo userrepo.Repository,
	client *asynq.Client,
	fixCostRepo fcrepo.Repository,
	txRepo txrepo.Repository,

) *Processor {
	return &Processor{
		txService:   txService,
		storage:     st,
		draftRepo:   dr,
		userRepo:    userRepo,
		client:      client,
		fixCostRepo: fixCostRepo,
		txRepo:      txRepo,
	}
}

func (p *Processor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(tasks.TaskBuildTransaction, p.handleBuildTransactionTask)
	mux.HandleFunc(tasks.TaskPurgeUser, p.HandlePurgeUser)
	mux.HandleFunc(tasks.TypeProcessFixCost, p.HandleFixCost)
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

	// IDEMPOTENT GUARD
	if exDraft.Status == m.DraftStatusWaitingConfirm {
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
					m.DraftStatusFailed,
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
	if exDraft.Status != m.DraftStatusProcessing {
		if err := p.draftRepo.UpdateStatus(
			ctx,
			payload.DraftID,
			m.DraftStatusProcessing,
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
			m.DraftStatusFailed,
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
	exDraft.Status = m.DraftStatusWaitingConfirm
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

func (p *Processor) HandleFixCost(ctx context.Context, t *asynq.Task) error {
	var payload tasks.ProcessFixCostPayload

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf(
		"[FIXCOST] start id=%s user=%s date=%s",
		payload.FixCostID,
		payload.UserID,
		payload.Date,
	)

	return p.processOneByID(
		ctx,
		payload.FixCostID,
		payload.Date,
		payload.UserID,
	)
}

func (p *Processor) processOne(ctx context.Context, fc *e.FixCost) error {

	_, err := p.txService.CreateFromFixCost(ctx, fc)
	if err != nil {
		log.Printf(
			"[FIXCOST] create tx failed id=%s err=%v",
			fc.FixCostID,
			err,
		)
		return err
	}

	fc.RunCount += 1
	next := fcsvc.CalculateNextRun(*fc)

	if fc.EndDate != nil {
		nextTruncate := next.Truncate(24 * time.Hour)
		endTruncate := fc.EndDate.Truncate(24 * time.Hour)
		if nextTruncate.After(endTruncate) {
			fc.Status = "finished"
			fc.LastRunAt = &fc.NextRunDate
			return p.fixCostRepo.Update(ctx, fc)
		}
	}

	if fc.MaxRun != nil {
		if fc.RunCount >= *fc.MaxRun {
			fc.Status = "finished"
			fc.LastRunAt = &fc.NextRunDate
			return p.fixCostRepo.Update(ctx, fc)
		}
	}

	fc.LastRunAt = &fc.NextRunDate
	fc.NextRunDate = next

	log.Printf(
		"[FIXCOST] success id=%s user=%s next_run=%s",
		fc.FixCostID,
		fc.UserID,
		fc.NextRunDate,
	)

	return p.fixCostRepo.Update(ctx, fc)
}

func (p *Processor) processOneByID(
	ctx context.Context,
	id string,
	date time.Time,
	userId string,
) error {

	fc, err := p.fixCostRepo.FindByID(ctx, id, userId)
	if err != nil {
		return err
	}

	loc := time.UTC
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	nextRun := date.In(loc).Truncate(24 * time.Hour)

	log.Printf("rundate: %s, today: %s", nextRun, today)

	// หยุดถ้ายังไม่ถึงวัน
	if nextRun.After(today) {
		return nil
	}

	// กันซ้ำ
	tx, err := p.txRepo.FindByFixCostIDAndRunDate(
		&m.TransactionFixCostSearchParams{
			FixCostID: fc.FixCostID,
			RunDate:   nextRun,
			UserID:    fc.UserID,
		},
	)
	if err != nil {
		return err
	}

	if tx == nil {
		if err := p.processOne(ctx, fc); err != nil {
			return err
		}
	} else {
		// ถ้ามี tx แล้ว → ขยับ next อย่างเดียว
		fc.LastRunAt = &fc.NextRunDate
		fc.NextRunDate = fcsvc.CalculateNextRun(*fc)

		if err := p.fixCostRepo.Update(ctx, fc); err != nil {
			return err
		}
	}
	
	return nil
}

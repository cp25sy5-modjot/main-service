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
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	fcsvc "github.com/cp25sy5-modjot/main-service/internal/fix_cost/service"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
)

type Processor struct {
	txService txsvc.Service
	storage   storage.Storage

	draftRepo   *d.DraftRepository
	userRepo    *userrepo.Repository
	fixCostRepo *fcrepo.Repository
	client      *asynq.Client
}

func NewProcessor(
	txService txsvc.Service,
	st storage.Storage,
	dr *d.DraftRepository,
	userRepo *userrepo.Repository,
	fixCostRepo *fcrepo.Repository,
	client *asynq.Client,
) *Processor {
	return &Processor{
		txService:   txService,
		storage:     st,
		draftRepo:   dr,
		userRepo:    userRepo,
		fixCostRepo: fixCostRepo,
		client:      client,
	}
}

func (p *Processor) Register(mux *asynq.ServeMux) {
	mux.HandleFunc(tasks.TaskBuildTransaction, p.handleBuildTransactionTask)
	mux.HandleFunc(tasks.TaskPurgeUser, p.HandlePurgeUser)
	mux.HandleFunc(tasks.TaskRunFixCost, p.HandleRunFixCost)
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

func sameUTCDate(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

func (p *Processor) HandleRunFixCost(ctx context.Context, t *asynq.Task) error {
	var payload tasks.RunFixCostPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	log.Printf("[JOB fix_cost] run fix_cost_id=%s", payload.FixCostID)

	fc, err := p.fixCostRepo.FindByID(ctx, payload.FixCostID, payload.UserID)
	if err != nil {
		return err
	}

	// skip if not active
	if fc.Status != "active" {
		return nil
	}

	now := time.Now().UTC()

	// ❗ 1. ยังไม่ถึงวัน (เทียบแค่ date)
	if now.Before(fc.NextRunDate) && !sameUTCDate(now, fc.NextRunDate) {
		return nil
	}

	// ❗ 2. กัน run ซ้ำในวันเดียวกัน
	if fc.LastRunAt != nil && sameUTCDate(*fc.LastRunAt, fc.NextRunDate) {
		return nil
	}

	// ✅ 3. create transaction (ควรมี idempotency ใน DB ด้วย)
	_, err = p.txService.CreateFromFixCost(ctx, fc)
	if err != nil {
		return err
	}

	// ✅ 4. calculate next run
	next := fcsvc.CalculateNextRun(*fc)

	// ❗ 5. check end date
	if fc.EndDate != nil && next.After(*fc.EndDate) {
		log.Printf("[JOB fix_cost] reached end date → stop")

		// update last run เพื่อกันซ้ำ
		fc.LastRunAt = &fc.NextRunDate
		return p.fixCostRepo.Update(ctx, fc)
	}

	// ❗ 6. handle remaining runs
	if fc.RemainingRuns != nil {
		*fc.RemainingRuns = *fc.RemainingRuns - 1

		if *fc.RemainingRuns <= 0 {
			log.Printf("[JOB fix_cost] no remaining runs → stop")

			fc.LastRunAt = &fc.NextRunDate
			return p.fixCostRepo.Update(ctx, fc)
		}
	}

	// ✅ 7. update state
	fc.LastRunAt = &fc.NextRunDate
	fc.NextRunDate = next

	err = p.fixCostRepo.Update(ctx, fc)
	if err != nil {
		return err
	}

	// ✅ 8. schedule next job
	task, err := tasks.NewRunFixCostTask(fc.FixCostID, fc.UserID)
	if err != nil {
		return err
	}

	_, err = p.client.Enqueue(
		task,
		asynq.ProcessAt(next),
		asynq.TaskID("fixcost:"+fc.FixCostID),
		asynq.Unique(24*time.Hour), // 🔥 กันซ้ำ
	)

	return err
}

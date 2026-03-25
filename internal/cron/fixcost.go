package cron

import (
	"context"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	fcsvc "github.com/cp25sy5-modjot/main-service/internal/fix_cost/service"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
)

type FixCostProcessor struct {
	repo      *fcrepo.Repository
	fcsvc     fcsvc.Service
	txService txsvc.Service
}

func (p *FixCostProcessor) Run() {
	ctx := context.Background()

	fcs, _ := p.repo.FindDueFixCosts(ctx)

	for _, fc := range fcs {
		_ = p.processOne(ctx, fc)
	}
}

func (p *FixCostProcessor) processOne(ctx context.Context, fc *e.FixCost) error {

	if fc.LastRunAt != nil && sameUTCDate(*fc.LastRunAt, fc.NextRunDate) {
		return nil
	}

	_, err := p.txService.CreateFromFixCost(ctx, fc)
	if err != nil {
		return err
	}

	next := fcsvc.CalculateNextRun(*fc)

	if fc.EndDate != nil && next.After(*fc.EndDate) {
		fc.Status = "finished"
		fc.LastRunAt = &fc.NextRunDate
		return p.repo.Update(ctx, fc)
	}

	if fc.RemainingRuns != nil {
		*fc.RemainingRuns--

		if *fc.RemainingRuns <= 0 {
			fc.Status = "finished"
			fc.LastRunAt = &fc.NextRunDate
			return p.repo.Update(ctx, fc)
		}
	}

	// ✅ update state
	fc.LastRunAt = &fc.NextRunDate
	fc.NextRunDate = next

	return p.repo.Update(ctx, fc)
}

func sameUTCDate(a, b time.Time) bool {
	ay, am, ad := a.UTC().Date()
	by, bm, bd := b.UTC().Date()
	return ay == by && am == bm && ad == bd
}

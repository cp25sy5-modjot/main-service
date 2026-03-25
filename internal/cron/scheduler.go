package cron

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/robfig/cron/v3"

	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
)

type Scheduler struct {
	client *asynq.Client
	repo   *fcrepo.Repository
}

func NewScheduler(client *asynq.Client, repo *fcrepo.Repository) *Scheduler {
	return &Scheduler{client: client, repo: repo}
}
func (s *Scheduler) Start() {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		ctx := context.Background()

		fcs, err := s.repo.FindDueFixCosts(ctx)
		if err != nil {
			log.Printf("find due error: %v", err)
			return
		}

		for _, fc := range fcs {
			task, err := tasks.NewProcessFixCostTask(
				fc.FixCostID,
				fc.NextRunDate,
				fc.UserID,
			)
			if err != nil {
				continue
			}

			_, err = s.client.Enqueue(
				task,
				asynq.Unique(24*time.Hour),
			)
			if err != nil {
				log.Printf("enqueue error: %v", err)
			}
		}
	})

	c.Start()
}

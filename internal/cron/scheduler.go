package cron

import "github.com/robfig/cron/v3"

type Scheduler struct {
	fixCostProcessor *FixCostProcessor
}

func (s *Scheduler) Start() {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		s.fixCostProcessor.Run()
	})

	c.Start()
}


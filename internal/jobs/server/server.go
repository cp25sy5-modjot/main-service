package server

import (
	"log"

	"github.com/hibiken/asynq"
)

func NewAsynqClient(redisAddr string) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})
}

func NewAsynqServer(redisAddr string, concurrency int) *asynq.Server {
	if concurrency <= 0 {
		concurrency = 5
	}
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: concurrency},
	)
}

func RunServer(srv *asynq.Server, mux *asynq.ServeMux) {
	if err := srv.Run(mux); err != nil {
		log.Fatalf("asynq server error: %v", err)
	}
}

package main

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// internal
	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	"github.com/cp25sy5-modjot/main-service/internal/cron"
	"github.com/cp25sy5-modjot/main-service/internal/database"
	d "github.com/cp25sy5-modjot/main-service/internal/draft"
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/processor"
	jobsserver "github.com/cp25sy5-modjot/main-service/internal/jobs/server"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	"github.com/cp25sy5-modjot/main-service/internal/storage/localfs"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"

	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
)

func main() {
	conf := config.LoadConfig()

	// =========================
	// DATABASE
	// =========================
	db := database.NewPostgresDatabase(conf)

	// =========================
	// gRPC (AI SERVICE)
	// =========================
	grpcConn, err := grpc.NewClient(
		conf.AIService.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close()

	aiClient := pb.NewAiWrapperServiceClient(grpcConn)

	// =========================
	// REDIS
	// =========================
	redisAddr := "localhost:6379"
	if conf.Redis != nil && conf.Redis.Addr != "" {
		redisAddr = conf.Redis.Addr
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	// =========================
	// REPOSITORIES
	// =========================
	txRepo := txrepo.NewRepository(db.GetDb())
	txiRepo := txirepo.NewRepository(db.GetDb())
	catRepo := catrepo.NewRepository(db.GetDb())
	userRepo := userrepo.NewRepository(db.GetDb())
	fcRepo := fcrepo.NewRepository(db.GetDb())

	// =========================
	// SERVICES
	// =========================
	txService := txsvc.NewService(db.GetDb(), txRepo, txiRepo, catRepo, aiClient)

	// =========================
	// STORAGE
	// =========================
	uploadDir := conf.Storage.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	st, err := localfs.NewLocalStorage(uploadDir)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// =========================
	// DRAFT
	// =========================
	draftRepo := d.NewDraftRepository(rdb)

	// =========================
	// CRON (enqueue job)
	// =========================
	scheduler := cron.NewScheduler(asynqClient, fcRepo)
	scheduler.Start()

	// =========================
	// ASYNQ SERVER
	// =========================
	srv := jobsserver.NewAsynqServer(redisAddr, 5)
	mux := asynq.NewServeMux()

	// =========================
	// PROCESSOR (handlers)
	// =========================
	p := processor.NewProcessor(
		txService,
		st,
		draftRepo,
		userRepo,
		asynqClient,
		fcRepo,
		txRepo,
	)

	p.Register(mux)

	// =========================
	// START WORKER
	// =========================
	log.Printf("Starting worker with Redis at %s", redisAddr)
	jobsserver.RunServer(srv, mux)
}

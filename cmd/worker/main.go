package main

import (
	"log"

	"github.com/hibiken/asynq"

	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	"github.com/cp25sy5-modjot/main-service/internal/database"
	d "github.com/cp25sy5-modjot/main-service/internal/draft"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/processor"
	jobsserver "github.com/cp25sy5-modjot/main-service/internal/jobs/server"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	"github.com/cp25sy5-modjot/main-service/internal/storage/localfs"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conf := config.LoadConfig()

	// DB
	db := database.NewPostgresDatabase(conf)

	// gRPC AI client (same as API server)
	grpcConn, err := grpc.NewClient(
		conf.AIService.Url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}

	defer func() {
		if err := grpcConn.Close(); err != nil {
			log.Printf("failed to close grpc connection: %v", err)
		}
	}()

	aiClient := pb.NewAiWrapperServiceClient(grpcConn)

	// Services
	txRepo := txrepo.NewRepository(db.GetDb())
	txiRepo := txirepo.NewRepository(db.GetDb())
	catRepo := catrepo.NewRepository(db.GetDb())
	userRepo := userrepo.NewRepository(db.GetDb())

	txService := txsvc.NewService(db.GetDb(), txRepo, txiRepo, catRepo, aiClient)

	// Storage
	uploadDir := conf.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	st, err := localfs.NewLocalStorage(uploadDir)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	// Redis addr
	redisAddr := ""
	if conf.Redis != nil {
		redisAddr = conf.Redis.Addr
	}
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	// ===== REDIS FOR DRAFT =====
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// สร้าง draft repo
	draftRepo := d.NewDraftRepository(rdb)

	// Asynq server
	srv := jobsserver.NewAsynqServer(redisAddr, 5)
	mux := asynq.NewServeMux()

	// Job processor
	p := processor.NewProcessor(txService, st, draftRepo, userRepo)
	p.Register(mux)

	log.Printf("Starting worker with Redis at %s", redisAddr)
	jobsserver.RunServer(srv, mux)
}

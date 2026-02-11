package httpapi

import (
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/database"
	globalhandler "github.com/cp25sy5-modjot/main-service/internal/global_handler"
	"github.com/cp25sy5-modjot/main-service/internal/middleware"
	"github.com/cp25sy5-modjot/main-service/internal/shared/config"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/cp25sy5-modjot/main-service/internal/storage"
	"github.com/cp25sy5-modjot/main-service/internal/storage/localfs"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"

	// "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Start()
}

type fiberServer struct {
	app         *fiber.App
	db          database.Database
	rdb         *redis.Client
	conf        *config.Config
	aiClient    pb.AiWrapperServiceClient
	asynqClient *asynq.Client
	storage     storage.Storage
}

func NewFiberServer(conf *config.Config, db database.Database, aiClient pb.AiWrapperServiceClient) Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: globalhandler.GlobalErrorHandler,
	})
	initMiddleware(app)

	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: conf.Redis.Addr,
	})

	rdb := redis.NewClient(&redis.Options{
		Addr: conf.Redis.Addr,
	})

	uploadDir := conf.Storage.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	st, err := localfs.NewLocalStorage(uploadDir)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	return &fiberServer{
		app:         app,
		db:          db,
		rdb:         rdb,
		conf:        conf,
		aiClient:    aiClient,
		asynqClient: asynqClient,
		storage:     st,
	}
}

func (s *fiberServer) Start() {
	RegisterRoutes(s)

	url, _ := utils.AppUrlBuilder(s.conf)
	log.Printf("🚀 Server running on %s", url)
	log.Fatal(s.app.Listen(":" + s.conf.App.Port))
}

func initMiddleware(app *fiber.App) {
	app.Use(middleware.RequestIDMiddleware)
	app.Use(middleware.LoggerMiddleware)
	app.Use(middleware.EnforceUTC())
	// app.Use(swagger.New(swagger.ConfigDefault))
}

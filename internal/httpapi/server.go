package httpapi

import (
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/config"
	"github.com/cp25sy5-modjot/main-service/internal/database"
	"github.com/cp25sy5-modjot/main-service/internal/globalHandler"
	"github.com/cp25sy5-modjot/main-service/internal/middleware"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"

	// "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Start()
}

type fiberServer struct {
	app  *fiber.App
	db   database.Database
	conf *config.Config
	aiClient   pb.AiWrapperServiceClient
}

func NewFiberServer(conf *config.Config, db database.Database, aiClient pb.AiWrapperServiceClient) Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: globalHandler.GlobalErrorHandler,
	})

	// Middlewares
	initMiddleware(app)
	return &fiberServer{
		app:  app,
		db:   db,
		conf: conf,
		aiClient:   aiClient,
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
	// app.Use(swagger.New(swagger.ConfigDefault))
}

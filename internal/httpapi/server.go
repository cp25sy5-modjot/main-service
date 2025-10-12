package httpapi

import (
	"log"
	"modjot/internal/config"
	"modjot/internal/database"
	"modjot/internal/middleware"
	"modjot/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Server interface {
	Start()
}

type fiberServer struct {
	app  *fiber.App
	db   database.Database
	conf *config.Config
}

func NewFiberServer(conf *config.Config, db database.Database) Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.GlobalErrorHandler,
	})

	app.Use(middleware.RequestIDMiddleware)
	app.Use(middleware.LoggerMiddleware)

	return &fiberServer{
		app:  app,
		db:   db,
		conf: conf,
	}
}

func (s *fiberServer) Start() {
	RegisterRoutes(s)

	url, _ := utils.AppUrlBuilder(s.conf)
	log.Printf("🚀 Server running on %s", url)
	log.Fatal(s.app.Listen(":" + s.conf.App.Port))
}

package server

import (
	"log"
	"modjot/config"
	"modjot/database"
	"modjot/modules/receipt"
	"modjot/utils"

	"github.com/gofiber/fiber/v2"
)

type fiberServer struct {
	app  *fiber.App
	db   database.Database
	conf *config.Config
}

func NewFiberServer(conf *config.Config, db database.Database) Server {
	fiberApp := fiber.New()

	return &fiberServer{
		app:  fiberApp,
		db:   db,
		conf: conf,
	}
}

func (s *fiberServer) Start() {
	initializeReceiptHttpHandler(s)
	initializeHealthCheck(s)

	url, _ := utils.AppUrlBuilder(s.conf)
	log.Printf("🚀 Server running on %s", url)
	log.Fatal(s.app.Listen(":" + s.conf.App.Port))
}

func initializeReceiptHttpHandler(s *fiberServer) {
	// Initialize all layers
	receiptRepo := receipt.NewRepositoryPg(s.db.GetDb())
	receiptUsecase := receipt.NewUsecase(receiptRepo)
	receiptHandler := receipt.NewHandler(receiptUsecase)

	// Register routes
	api := s.app.Group("/api")
	api.Post("/receipts", receiptHandler.Create)
	api.Get("/receipts", receiptHandler.GetAll)
	api.Get("/receipts/:id", receiptHandler.GetByID)
	api.Put("/receipts/:id", receiptHandler.Update)
	api.Delete("/receipts/:id", receiptHandler.Delete)
}

func initializeHealthCheck(s *fiberServer) {
	s.app.Get("v1/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
}

package routes

import (
	"modjot/internal/receipt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB) {
	// Receipt
	receiptRepo := receipt.NewRepositoryPg(db)
	receiptService := receipt.NewService(receiptRepo)
	receiptHandler := receipt.NewHandler(receiptService)

	api := app.Group("/api")
	api.Post("/receipts", receiptHandler.Create)
	api.Get("/receipts", receiptHandler.GetAll)
	api.Get("/receipts/:id", receiptHandler.GetByID)
	api.Put("/receipts/:id", receiptHandler.Update)
	api.Delete("/receipts/:id", receiptHandler.Delete)
}

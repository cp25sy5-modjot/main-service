package httpapi

import (
	r "modjot/internal/response"
	"modjot/internal/transaction"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	s *fiberServer,
) {
	initializeHealthCheck(s)
	initializeTransactionRoutes(s)
}

func initializeTransactionRoutes(s *fiberServer) {
	// Initialize all layers
	transactionRepo := transaction.NewRepository(s.db.GetDb())
	transactionService := transaction.NewService(transactionRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	// Register routes
	api := s.app.Group("")
	api.Post("/transactions", transactionHandler.Create)
	api.Get("/transactions", transactionHandler.GetAll)
	api.Get("/transactions/:id", transactionHandler.GetByID)
	api.Put("/transactions/:id", transactionHandler.Update)
	api.Delete("/transactions/:id", transactionHandler.Delete)
}

func initializeHealthCheck(s *fiberServer) {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Health check passed")
	})
}

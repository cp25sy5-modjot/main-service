package httpapi

import (
	"modjot/internal/auth"
	r "modjot/internal/response"
	"modjot/internal/transaction"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	s *fiberServer,
) {
	initializeHealthCheck(s)
	initializeTransactionRoutes(s)
	initializeAuthRoutes(s)
}

func initializeTransactionRoutes(s *fiberServer) {
	// Initialize all layers
	transactionRepo := transaction.NewRepository(s.db.GetDb())
	transactionService := transaction.NewService(transactionRepo)
	transactionHandler := transaction.NewHandler(transactionService)

	// Register routes
	api := s.app.Group("/transactions")
	api.Use(auth.Protected(s.conf.Auth.AccessTokenSecret))

	api.Post("", transactionHandler.Create)
	api.Get("", transactionHandler.GetAll)
	api.Get("/:transaction_id/:product_id", transactionHandler.GetByID)
	api.Put("/:transaction_id/:product_id", transactionHandler.Update)
	api.Delete("/:transaction_id/:product_id", transactionHandler.Delete)
}

func initializeHealthCheck(s *fiberServer) {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Health check passed")
	})
}

func initializeAuthRoutes(s *fiberServer) {
	// Initialize all layers
	// authRepo := auth.NewRepository(s.db.GetDb())
	// authService := auth.NewService(authRepo, s.conf)
	// authHandler := auth.NewHandler(authService)

	// Register routes
	api := s.app.Group("/auth")
	// Mock login route for testing purposes
	api.Post("/login", func(c *fiber.Ctx) error {
		return auth.LoginHandler(c, s.conf.Auth)
	})
	api.Post("/refresh", func(c *fiber.Ctx) error {
		return auth.RefreshHandler(c, s.conf.Auth)
	})
	// Google OAuth routes (placeholders)
	api.Get("/google", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Google login endpoint")
	})
	// api.Post("/google/callback", authHandler.GoogleCallback)
}

package httpapi

import (
	"modjot/internal/auth"
	r "modjot/internal/response"
	"modjot/internal/transaction"
	"modjot/internal/user"

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
	api := s.app.Group("/v1/transactions")
	api.Use(auth.Protected(s.conf.Auth.AccessTokenSecret))

	api.Post("/manual", transactionHandler.Create)
	api.Get("", transactionHandler.GetAll)
	api.Get("/transaction/:transaction_id/product/:product_id", transactionHandler.GetByID)
	api.Put("/transaction/:transaction_id/product/:product_id", transactionHandler.Update)
	api.Delete("/transaction/:transaction_id/product/:product_id", transactionHandler.Delete)
}

func initializeHealthCheck(s *fiberServer) {
	s.app.Get("/health", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Health check passed")
	})
}

func initializeAuthRoutes(s *fiberServer) {
	userRepo := user.NewRepository(s.db.GetDb())
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Register user routes
	userApi := s.app.Group("/v1/user")
	userApi.Use(auth.Protected(s.conf.Auth.AccessTokenSecret))
	userApi.Put("/:id", userHandler.Update)
	userApi.Delete("/:id", userHandler.Delete)

	authApi := s.app.Group("/v1/auth")
	authApi.Post("/mock-login", func(c *fiber.Ctx) error {
		return auth.MockLoginHandler(c, s.conf.Auth)
	})
	authApi.Post("/refresh-token", func(c *fiber.Ctx) error {
		return auth.RefreshHandler(c, s.conf.Auth)
	})
	authApi.Get("/google", func(c *fiber.Ctx) error {
		return auth.HandleGoogleTokenExchange(c,userService, s.conf)
	})

}

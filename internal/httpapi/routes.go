package httpapi

import (
	"context"
	"log"
	"time"

	"github.com/cp25sy5-modjot/main-service/internal/auth"
	catHandler "github.com/cp25sy5-modjot/main-service/internal/category/handler"
	catRepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"

	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	r "github.com/cp25sy5-modjot/main-service/internal/response/success"
	tranHandler "github.com/cp25sy5-modjot/main-service/internal/transaction/handler"
	tranRepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	tranSvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	userHandler "github.com/cp25sy5-modjot/main-service/internal/user/handler"
	userRepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	userSvc "github.com/cp25sy5-modjot/main-service/internal/user/service"

	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"

	"github.com/gofiber/fiber/v2"
)

type Services struct {
	UserService        *userSvc.Service
	TransactionService *tranSvc.Service
	CategoryService    *catSvc.Service
}

func RegisterRoutes(
	s *fiberServer,
) {
	services := initializeServices(s)

	initializeHealthCheckRoutes(s)
	initializeTransactionRoutes(s, services)
	initializeAuthRoutes(s, services)
	initializeCategoryRoutes(s, services)
}
func initializeServices(s *fiberServer) *Services {

	// Category Service
	categoryRepo := catRepo.NewRepository(s.db.GetDb())
	categorySvc := catSvc.NewService(categoryRepo)

	// User Service
	userRepo := userRepo.NewRepository(s.db.GetDb())
	userSvc := userSvc.NewService(userRepo, categorySvc)
	
	// Transaction Service
	transactionRepo := tranRepo.NewRepository(s.db.GetDb())
	transactionSvc := tranSvc.NewService(transactionRepo, categorySvc, s.aiClient)

	return &Services{
		UserService:        userSvc,
		TransactionService: transactionSvc,
		CategoryService:    categorySvc,
	}
}

func initializeHealthCheckRoutes(s *fiberServer) {
	s.app.Get("/v1/health", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Health check passed")
	})
	s.app.Get("/v1/health/grpc", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 sec timeout for upload
		defer cancel()
		log.Println("Performing gRPC health check...")
		resp, err := s.aiClient.Check(ctx, &pb.HealthCheckRequest{Name: "main-service"})
		if err != nil {
			return err
		}
		return r.OK(c, resp, "gRPC health check passed")
	})
}

func initializeAuthRoutes(s *fiberServer, services *Services) {
	userHandler := userHandler.NewHandler(services.UserService)

	// Register user routes
	userApi := s.app.Group("/v1/user")
	userApi.Use(jwt.Protected(s.conf.Auth.AccessTokenSecret))

	userApi.Get("", userHandler.GetSelf)
	userApi.Put("", userHandler.Update)
	userApi.Delete("", userHandler.Delete)

	authApi := s.app.Group("/v1/auth")
	authApi.Post("/mock-login", func(c *fiber.Ctx) error {
		return auth.MockLoginHandler(c, services.UserService, s.conf.Auth)
	})
	authApi.Post("/refresh-token", func(c *fiber.Ctx) error {
		return auth.RefreshHandler(c, s.conf.Auth)
	})
	authApi.Post("/google", func(c *fiber.Ctx) error {
		return auth.HandleGoogleTokenExchange(c, services.UserService, s.conf)
	})

}

func initializeTransactionRoutes(s *fiberServer, services *Services) {
	transactionHandler := tranHandler.NewHandler(services.TransactionService)

	// Register routes
	txApi := s.app.Group("/v1/transaction")
	txApi.Use(jwt.Protected(s.conf.Auth.AccessTokenSecret))

	txApi.Post("/manual", transactionHandler.Create)
	txApi.Post("/upload", transactionHandler.UploadImage)
	txApi.Get("", transactionHandler.GetAll)
	txApi.Get("/:transaction_id/product/:product_id", transactionHandler.GetByID)
	txApi.Put("/:transaction_id/product/:product_id", transactionHandler.Update)
	txApi.Delete("/:transaction_id/product/:product_id", transactionHandler.Delete)
}

func initializeCategoryRoutes(s *fiberServer, services *Services) {
	categoryHandler := catHandler.NewHandler(services.CategoryService)

	// Register routes
	api := s.app.Group("/v1/category")
	api.Use(jwt.Protected(s.conf.Auth.AccessTokenSecret))

	api.Post("", categoryHandler.Create)
	api.Get("", categoryHandler.GetAll)
	api.Put("/:id", categoryHandler.Update)
	api.Delete("/:id", categoryHandler.Delete)
}

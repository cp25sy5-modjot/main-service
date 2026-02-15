package httpapi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cp25sy5-modjot/main-service/internal/auth"
	cathandler "github.com/cp25sy5-modjot/main-service/internal/category/handler"
	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	catsvc "github.com/cp25sy5-modjot/main-service/internal/category/service"

	draft "github.com/cp25sy5-modjot/main-service/internal/draft"
	favhandler "github.com/cp25sy5-modjot/main-service/internal/favorite_item/handler"
	favrepo "github.com/cp25sy5-modjot/main-service/internal/favorite_item/repository"
	favsvc "github.com/cp25sy5-modjot/main-service/internal/favorite_item/service"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	overviewhandler "github.com/cp25sy5-modjot/main-service/internal/overview/handler"
	overviewrepo "github.com/cp25sy5-modjot/main-service/internal/overview/repository"
	overviewsvc "github.com/cp25sy5-modjot/main-service/internal/overview/service"
	r "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	txhandler "github.com/cp25sy5-modjot/main-service/internal/transaction/handler"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txsvc "github.com/cp25sy5-modjot/main-service/internal/transaction/service"
	txihandler "github.com/cp25sy5-modjot/main-service/internal/transaction_item/handler"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
	txisvc "github.com/cp25sy5-modjot/main-service/internal/transaction_item/service"
	userhandler "github.com/cp25sy5-modjot/main-service/internal/user/handler"
	userepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	usersvc "github.com/cp25sy5-modjot/main-service/internal/user/service"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"github.com/gofiber/fiber/v2"
)

type Services struct {
	UserService            usersvc.Service
	TransactionService     txsvc.Service
	TransactionItemService txisvc.Service
	CategoryService        catsvc.Service

	OverviewService overviewsvc.Service
	DraftService    draft.Service
	FavoriteService favsvc.Service
}

type Repositories struct {
	UserRepo            *userepo.Repository
	TransactionRepo     *txrepo.Repository
	TransactionItemRepo *txirepo.Repository
	CategoryRepo        *catrepo.Repository

	OverviewRepo *overviewrepo.Repository
	DraftRepo    *draft.DraftRepository
	FavoriteRepo *favrepo.Repository
}

func RegisterRoutes(s *fiberServer) {

	repositories := initializeRepositories(s)
	services := initializeServices(s, repositories)

	// ⭐ สร้างครั้งเดียวตรงนี้
	authMiddleware := jwt.Protected(
		s.conf.Auth.AccessTokenSecret,
		repositories.UserRepo,
	)

	initializeHealthCheckRoutes(s)
	initializeFileRoutes(s)

	initializeDraftRoutes(s, services, authMiddleware)
	initializeTransactionRoutes(s, services, authMiddleware)
	initializeTransactionItemRoutes(s, services, authMiddleware)
	initializeAuthRoutes(s, services, authMiddleware)
	initializeCategoryRoutes(s, services, authMiddleware)
	initializeOverviewRoutes(s, services, authMiddleware)
	initializeFavoriteRoutes(s, services, authMiddleware)
}

func initializeRepositories(s *fiberServer) *Repositories {
	categoryRepo := catrepo.NewRepository(s.db.GetDb())
	userRepo := userepo.NewRepository(s.db.GetDb())
	transactionRepo := txrepo.NewRepository(s.db.GetDb())
	transactionItemRepo := txirepo.NewRepository(s.db.GetDb())
	overviewRepo := overviewrepo.NewRepository(s.db.GetDb())

	draftRepo := draft.NewDraftRepository(s.rdb)
	favRepo := favrepo.NewRepository(s.db.GetDb())

	return &Repositories{
		UserRepo:            userRepo,
		TransactionRepo:     transactionRepo,
		TransactionItemRepo: transactionItemRepo,
		CategoryRepo:        categoryRepo,

		OverviewRepo: overviewRepo,
		DraftRepo:    draftRepo,
		FavoriteRepo: favRepo,
	}
}

func initializeServices(s *fiberServer, repositories *Repositories) *Services {
	categoryRepo := repositories.CategoryRepo
	userRepo := repositories.UserRepo
	transactionRepo := repositories.TransactionRepo
	transactionItemRepo := repositories.TransactionItemRepo
	overviewRepo := repositories.OverviewRepo

	draftRepo := repositories.DraftRepo
	favRepo := repositories.FavoriteRepo

	transactionSvc := txsvc.NewService(
		s.db.GetDb(),
		transactionRepo,
		transactionItemRepo,
		categoryRepo,
		s.aiClient,
	)

	// 👇 สำคัญ: inject createInternal
	draftSvc := draft.NewService(
		draftRepo,
		categoryRepo,
		s.storage,
		s.conf.Storage.SignedURLSecret,
		transactionSvc.CreateInternal,
	)

	favSvc := favsvc.NewService(
		s.db.GetDb(),
		favRepo,
	)

	// ======================

	categorySvc := catsvc.NewService(categoryRepo, transactionRepo)
	userSvc := usersvc.NewService(userRepo, s.asynqClient)
	transactionItemSvc := txisvc.NewService(transactionItemRepo, transactionRepo)
	overviewSvc := overviewsvc.NewService(overviewRepo)

	return &Services{
		UserService:            userSvc,
		TransactionService:     transactionSvc,
		TransactionItemService: transactionItemSvc,
		CategoryService:        categorySvc,
		OverviewService:        overviewSvc,

		DraftService:    draftSvc,
		FavoriteService: favSvc,
	}
}

func initializeHealthCheckRoutes(s *fiberServer) {
	s.app.Get("/v1/health", func(c *fiber.Ctx) error {
		return r.OK(c, nil, "Health check passed")
	})
	s.app.Get("/v1/health/grpc", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 sec timeout for upload
		defer cancel()
		resp, err := s.aiClient.Check(ctx, &pb.HealthCheckRequest{Name: "main-service"})
		if err != nil {
			return err
		}
		return r.OK(c, resp, "gRPC health check passed")
	})
}

func initializeAuthRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	userHandler := userhandler.NewHandler(services.UserService)

	// Register user routes
	userApi := s.app.Group("/v1/user")
	userApi.Use(authMiddleware)

	userApi.Get("", userHandler.GetSelf)
	userApi.Put("", userHandler.Update)
	userApi.Delete("", userHandler.Delete)

	authApi := s.app.Group("/v1/auth")
	authApi.Post("/mock-login", func(c *fiber.Ctx) error {
		return auth.MockLoginHandler(c, services.UserService, services.CategoryService, s.conf.Auth)
	})
	authApi.Post("/mock-restore", func(c *fiber.Ctx) error {
		return auth.MockRestoreHandler(c, services.UserService, s.conf.Auth)
	})
	authApi.Post("/refresh-token", func(c *fiber.Ctx) error {
		return auth.RefreshHandler(c, services.UserService, s.conf.Auth)
	})
	authApi.Post("/google", func(c *fiber.Ctx) error {
		return auth.HandleGoogleTokenExchange(c, services.UserService, services.CategoryService, s.conf)
	})
	authApi.Post("/google/restore", func(c *fiber.Ctx) error {
		return auth.HandleRestore(c, services.UserService, s.conf)
	})

}

func initializeTransactionRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	transactionHandler := txhandler.NewHandler(
		services.TransactionService,
		s.asynqClient,
		s.storage,
		services.DraftService,
		services.FavoriteService,
	)

	txApi := s.app.Group("/v1/transaction")
	txApi.Use(authMiddleware)

	txApi.Post("/manual", transactionHandler.Create)
	txApi.Post("/upload", transactionHandler.UploadImage) // async
	txApi.Get("", transactionHandler.GetAll)
	txApi.Get("/:transaction_id", transactionHandler.GetByID)
	txApi.Patch("/:transaction_id", transactionHandler.Update)
	txApi.Delete("/:transaction_id", transactionHandler.Delete)
}

func initializeTransactionItemRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	transactionItemHandler := txihandler.NewHandler(services.TransactionItemService)

	// Register routes
	txItemApi := s.app.Group("/v1/transaction/:transaction_id/item")
	txItemApi.Use(authMiddleware)

	txItemApi.Get("/:item_id", transactionItemHandler.GetByID)
	txItemApi.Put("/:item_id", transactionItemHandler.Update)
	txItemApi.Delete("/:item_id", transactionItemHandler.Delete)
}

func initializeCategoryRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	categoryHandler := cathandler.NewHandler(services.CategoryService)

	// Register routes
	api := s.app.Group("/v1/category")
	api.Use(authMiddleware)

	api.Post("", categoryHandler.Create)
	api.Get("", categoryHandler.GetAll)
	api.Get("/:id", categoryHandler.GetByID)
	api.Put("/:id", categoryHandler.Update)
	api.Delete("/:id", categoryHandler.Delete)
}

func initializeOverviewRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	overviewHandler := overviewhandler.NewHandler(services.OverviewService)

	// Register routes
	api := s.app.Group("/v1/overview")
	api.Use(authMiddleware)

	api.Get("", overviewHandler.GetOverview)
}

func initializeDraftRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	handler := draft.NewHandler(
		services.DraftService,
	)

	api := s.app.Group("/v1/draft")
	api.Use(authMiddleware)

	api.Get("", handler.ListDraft)
	api.Get("/stats", handler.GetDraftStats)
	api.Get("/:draftID", handler.GetDraft)
	api.Post("/:draftID/confirm", handler.Confirm)
	api.Get("/:draftID/image-url", handler.GetDraftImageURL)
	api.Delete("/:draftID", handler.DeleteDraft)

}

func initializeFavoriteRoutes(s *fiberServer, services *Services, authMiddleware fiber.Handler) {
	favHandler := favhandler.NewHandler(services.FavoriteService)

	// Register routes
	api := s.app.Group("/v1/favorites")
	api.Use(authMiddleware)

	api.Post("", favHandler.Create)
	api.Get("", favHandler.GetAll)
	api.Get("/:id", favHandler.GetByID)
	api.Put("/:id", favHandler.Update)
	api.Delete("/:id", favHandler.Delete)
	api.Post("/reorder", favHandler.ReOrder)
}

func initializeFileRoutes(s *fiberServer) {

	secret := s.conf.Storage.SignedURLSecret
	baseDir := s.conf.Storage.UploadDir // ต้องมีใน config

	s.app.Get("/v1/files/*", func(c *fiber.Ctx) error {

		path := c.Params("*")
		expires := c.Query("expires")
		sig := c.Query("sig")

		if path == "" || expires == "" || sig == "" {
			return fiber.NewError(400, "invalid request")
		}

		// parse expiry
		expInt, err := strconv.ParseInt(expires, 10, 64)
		if err != nil {
			return fiber.NewError(400, "invalid expiry")
		}

		if time.Now().Unix() > expInt {
			return fiber.NewError(403, "url expired")
		}

		// validate signature
		expected := generateHMAC(
			fmt.Sprintf("%s:%d", path, expInt),
			secret,
		)

		if !hmac.Equal([]byte(sig), []byte(expected)) {
			return fiber.NewError(403, "invalid signature")
		}

		fullPath := filepath.Join(baseDir, path)

		return c.SendFile(fullPath)
	})
}

func generateHMAC(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

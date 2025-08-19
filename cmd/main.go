package main

import (
	"log"

	"modjot/config"
	"modjot/database"
	"modjot/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	db, err := database.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("DB migration failed:", err)
	}

	// Start Fiber
	app := fiber.New()

	// Register routes
	routes.Register(app, db)

	log.Printf("🚀 Server running on %s", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}

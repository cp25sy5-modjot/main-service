#!/bin/bash

# =========================================================
# Go Fiber + PostgreSQL Project Dependency Installer
# =========================================================
# This script will install all required Go libraries
# for the project (web framework, ORM, config, etc.)
# Run with:
#   chmod +x setup.sh
#   ./setup.sh
# =========================================================

echo "🚀 Initializing Go module..."
go mod init go-fiber-postgres-app

echo "📦 Installing Fiber (web framework)..."
# Fiber is a fast web framework (like Express.js in Node.js)
go get github.com/gofiber/fiber/v2

echo "📦 Installing Fiber middlewares (logger, recover, cors)..."
# logger  -> logs requests (method, path, latency)
# recover -> prevents crashes by recovering from panics
# cors    -> handles cross-origin requests
go get github.com/gofiber/fiber/v2/middleware/logger
go get github.com/gofiber/fiber/v2/middleware/recover
go get github.com/gofiber/fiber/v2/middleware/cors

echo "📦 Installing GORM (ORM for PostgreSQL)..."
# GORM is the ORM (Object Relational Mapper)
# postgres driver allows GORM to talk with PostgreSQL
go get gorm.io/gorm
go get gorm.io/driver/postgres

echo "📦 Installing godotenv..."
# godotenv allows loading configuration from .env files
go get github.com/joho/godotenv

echo "📦 Installing validator..."
# validator provides struct field validation (email, required, etc.)
go get github.com/go-playground/validator/v10

echo "📦 Installing UUID support..."
# UUID is useful if you prefer UUID instead of numeric IDs
go get github.com/google/uuid

echo "✅ All dependencies installed successfully!"

echo "📦 Tidying up go.mod & go.sum..."
go mod tidy

echo "🎉 Setup complete! Now you can run your app with:"
echo "    go run cmd/main.go"

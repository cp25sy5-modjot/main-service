package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Attach a unique trace ID to each request
func RequestIDMiddleware(c *fiber.Ctx) error {
	rid := c.Get("X-Request-ID")
	if rid == "" {
		rid = uuid.NewString()
	}
	c.Locals("request_id", rid)
	c.Set("X-Request-ID", rid)
	return c.Next()
}

// Simple console logger
func LoggerMiddleware(c *fiber.Ctx) error {
	log.Printf("[%s] %s %s", c.Locals("request_id"), c.Method(), c.OriginalURL())
	return c.Next()
}

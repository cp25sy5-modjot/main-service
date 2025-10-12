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

// Central error handler for all Fiber errors
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"
	typ := "internal_error"

	if fe, ok := err.(*fiber.Error); ok {
		code = fe.Code
		msg = fe.Message
		switch code {
		case fiber.StatusNotFound:
			typ = "not_found"
		case fiber.StatusUnauthorized:
			typ = "unauthorized"
		case fiber.StatusBadRequest:
			typ = "bad_request"
		}
	}

	return WriteError(c, code, msg, typ, err.Error(), nil)
}

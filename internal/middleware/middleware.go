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
// func GlobalErrorHandler(c *fiber.Ctx, err error) error {
// 	code := fiber.StatusInternalServerError

// 	if fe, ok := err.(*fiber.Error); ok {
// 		code = fe.Code
// 		switch code {
// 		case fiber.StatusNotFound:
// 			return r.NotFound(c, fe.Message)
// 		case fiber.StatusUnauthorized:
// 			return r.Unauthorized(c, fe.Message)
// 		case fiber.StatusBadRequest:
// 			return r.BadRequest(c, fe.Message)
// 		case fiber.StatusForbidden:
// 			return r.Forbidden(c, fe.Message)
// 		case fiber.StatusConflict:
// 			return r.Conflict(c, fe.Message)
// 		case fiber.StatusUnprocessableEntity:
// 			return r.UnprocessableEntity(c, fe.Message)
// 		case fiber.StatusTooManyRequests:
// 			return r.TooManyRequests(c, fe.Message)
// 		}

// 	}
// 	return r.InternalServerError(c, err.Error())
// }

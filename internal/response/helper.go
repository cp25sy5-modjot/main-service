package response

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// OK → 200
func OK(c *fiber.Ctx, data any, msg ...string) error {
	m := firstOrEmpty(msg, "OK")
	return WriteSuccess(c, fiber.StatusOK, data, m)
}

// Created → 201
func Created(c *fiber.Ctx, data any, msg ...string) error {
	m := firstOrEmpty(msg, "Created successfully")
	return WriteSuccess(c, fiber.StatusCreated, data, m)
}

// Accepted → 202
func Accepted(c *fiber.Ctx, data any, msg ...string) error {
	m := firstOrEmpty(msg, "Accepted")
	return WriteSuccess(c, fiber.StatusAccepted, data, m)
}

// NoContent → 204
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// OKWithMeta → 200 + meta (เช่น pagination)
func OKWithMeta(c *fiber.Ctx, data any, meta map[string]any, msg ...string) error {
	m := firstOrEmpty(msg, "OK")
	env := Response{
		Status:    "success",
		Code:      fiber.StatusOK,
		Message:   m,
		Data:      data,
		Meta:      meta,
		TraceID:   getTraceID(c),
		Timestamp: time.Now().UTC(),
	}
	return c.Status(fiber.StatusOK).JSON(env)
}

// BadRequest → 400
func BadRequest(c *fiber.Ctx, detail string, err error) error {
	return WriteError(c, fiber.StatusBadRequest, "Bad Request", "bad_request", detail, MapValidationErrors(err))
}

// Unauthorized → 401
func Unauthorized(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "unauthorized", detail, nil)
}

// Forbidden → 403
func Forbidden(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusForbidden, "Forbidden", "forbidden", detail, nil)
}

// NotFound → 404
func NotFound(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusNotFound, "Resource Not Found", "not_found", detail, nil)
}

// Conflict → 409 (เช่น unique constraint)
func Conflict(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusConflict, "Conflict", "conflict", detail, nil)
}

// UnprocessableEntity → 422 (validation error)
func UnprocessableEntity(c *fiber.Ctx, detail string, err error) error {
	return WriteError(c, fiber.StatusUnprocessableEntity, "Validation Failed", "validation_error", detail, MapValidationErrors(err))
}

// TooManyRequests -> 429 (rate limit)
func TooManyRequests(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusTooManyRequests, "Too Many Requests", "rate_limited", detail, nil)
}

// InternalServerError → 500
func InternalServerError(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", detail, nil)
}

// BadGateway → 502
func BadGateway(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusBadGateway, "Bad Gateway", "bad_gateway", detail, nil)
}

// ServiceUnavailable → 503
func ServiceUnavailable(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusServiceUnavailable, "Service Unavailable", "service_unavailable", detail, nil)
}

// GatewayTimeout → 504
func GatewayTimeout(c *fiber.Ctx, detail string) error {
	return WriteError(c, fiber.StatusGatewayTimeout, "Gateway Timeout", "gateway_timeout", detail, nil)
}

// Utility
func firstOrEmpty(msg []string, fallback string) string {
	if len(msg) > 0 && msg[0] != "" {
		return msg[0]
	}
	return fallback
}

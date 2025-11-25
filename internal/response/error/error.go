package errorresponse

import (
	"github.com/cp25sy5-modjot/main-service/internal/response"
	v "github.com/cp25sy5-modjot/main-service/internal/validator"
	"github.com/gofiber/fiber/v2"
)

// BadRequest → 400
func BadRequest(c *fiber.Ctx, detail string, err error) error {
	return response.WriteError(c, fiber.StatusBadRequest, "Bad Request", "bad_request", detail, v.MapValidationErrors(err))
}

// Unauthorized → 401
func Unauthorized(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusUnauthorized, "Unauthorized", "unauthorized", detail, nil)
}

// Forbidden → 403
func Forbidden(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusForbidden, "Forbidden", "forbidden", detail, nil)
}

// NotFound → 404
func NotFound(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusNotFound, "Resource Not Found", "not_found", detail, nil)
}

// Conflict → 409 (เช่น unique constraint)
func Conflict(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusConflict, "Conflict", "conflict", detail, nil)
}

// UnprocessableEntity → 422 (validation error)
func UnprocessableEntity(c *fiber.Ctx, detail string, err error) error {
	return response.WriteError(c, fiber.StatusUnprocessableEntity, "Validation Failed", "validation_error", detail, v.MapValidationErrors(err))
}

// TooManyRequests -> 429 (rate limit)
func TooManyRequests(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusTooManyRequests, "Too Many Requests", "rate_limited", detail, nil)
}

// InternalServerError → 500
func InternalServerError(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", detail, nil)
}

// BadGateway → 502
func BadGateway(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusBadGateway, "Bad Gateway", "bad_gateway", detail, nil)
}

// ServiceUnavailable → 503
func ServiceUnavailable(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusServiceUnavailable, "Service Unavailable", "service_unavailable", detail, nil)
}

// GatewayTimeout → 504
func GatewayTimeout(c *fiber.Ctx, detail string) error {
	return response.WriteError(c, fiber.StatusGatewayTimeout, "Gateway Timeout", "gateway_timeout", detail, nil)
}

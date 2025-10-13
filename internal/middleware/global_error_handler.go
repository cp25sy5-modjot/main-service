package middleware

import (
	"context"
	"errors"
	"log"
	"os"

	"database/sql"
	r "modjot/internal/response"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var isProd = os.Getenv("APP_ENV") == "prod"

func GlobalErrorHandler(c *fiber.Ctx, err error) error {

	// 1) fiber.Error
	var fe *fiber.Error
	if errors.As(err, &fe) {
		switch fe.Code {
		case fiber.StatusNotFound:
			return r.NotFound(c, fe.Message)
		case fiber.StatusUnauthorized:
			return r.Unauthorized(c, fe.Message)
		case fiber.StatusBadRequest:
			return r.BadRequest(c, fe.Message)
		case fiber.StatusForbidden:
			return r.Forbidden(c, fe.Message)
		case fiber.StatusConflict:
			return r.Conflict(c, fe.Message)
		case fiber.StatusUnprocessableEntity:
			return r.UnprocessableEntity(c, fe.Message, r.MapValidationErrors(err)...)
		case fiber.StatusTooManyRequests:
			return r.TooManyRequests(c, fe.Message)
		default:
			return r.WriteError(c, fe.Code, fe.Message, "http_error", safeDetail(err), nil)
		}
	}
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, gorm.ErrRecordNotFound) {
		return r.NotFound(c, "Resource Not Found")
	}

	if errors.Is(err, context.Canceled) {
		return r.WriteError(c, 499, "Client Closed Request", "client_closed", safeDetail(err), nil)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return r.GatewayTimeout(c, "Upstream timeout")
	}

	return r.InternalServerError(c, safeDetail(err))
}

func safeDetail(err error) string {
	if isProd {
		log.Println("detail not shown in production:")
		return "" // ปิด detail ใน prod
	}
	return err.Error() // โชว์ใน dev
}

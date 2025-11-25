package globalHandler

import (
	"errors"
	"log"
	"os"

	r "github.com/cp25sy5-modjot/main-service/internal/response"
	eresp "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	v "github.com/cp25sy5-modjot/main-service/internal/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var isProd = os.Getenv("APP_ENV") == "prod"

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	// log.Printf("Error caught by global handler: (%T) %v", err, err)

	var valErr *utils.ValidationError
	if errors.As(err, &valErr) {
		return r.WriteError(c, fiber.StatusUnprocessableEntity, valErr.Message, "validation_failed", valErr.Message, v.MapValidationErrors(valErr.OriginalErr))
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		switch fe.Code {
		case fiber.StatusBadRequest:
			return eresp.BadRequest(c, fe.Message, nil)
		case fiber.StatusNotFound:
			return eresp.NotFound(c, fe.Message)
		case fiber.StatusUnauthorized:
			return eresp.Unauthorized(c, fe.Message)
		case fiber.StatusForbidden:
			return eresp.Forbidden(c, fe.Message)
		case fiber.StatusConflict:
			return eresp.Conflict(c, fe.Message)
		case fiber.StatusTooManyRequests:
			return eresp.TooManyRequests(c, fe.Message)
		case fiber.StatusInternalServerError:
			return eresp.InternalServerError(c, fe.Message)
		case fiber.StatusBadGateway:
			return eresp.BadGateway(c, fe.Message)
		case fiber.StatusServiceUnavailable:
			return eresp.ServiceUnavailable(c, fe.Message)
		case fiber.StatusGatewayTimeout:
			return eresp.GatewayTimeout(c, fe.Message)
		default:
			log.Printf("Unhandled fiber error code: %d", fe.Code)
			return r.WriteError(c, fe.Code, "Error", "error", safeDetail(err), nil)
		}
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return eresp.NotFound(c, "Resource not found")
	}

	return eresp.InternalServerError(c, "An unexpected error occurred")
}

func safeDetail(err error) string {
	if isProd {
		log.Println("detail not shown in production:")
		return "" // ปิด detail ใน prod
	}
	return err.Error() // โชว์ใน dev
}

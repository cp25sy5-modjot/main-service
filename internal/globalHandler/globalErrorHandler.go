package globalHandler

import (
	"errors"
	"log"
	"os"

	r "github.com/cp25sy5-modjot/main-service/internal/response"
	errorResp "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	v "github.com/cp25sy5-modjot/main-service/internal/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var isProd = os.Getenv("APP_ENV") == "prod"

func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	log.Printf("Error caught by global handler: (%T) %v", err, err)

	var valErr *utils.ValidationError
	if errors.As(err, &valErr) {
		return r.WriteError(c, fiber.StatusUnprocessableEntity, valErr.Message, "validation_failed", valErr.Message, v.MapValidationErrors(valErr.OriginalErr))
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		switch fe.Code {
		case fiber.StatusBadRequest:
			return errorResp.BadRequest(c, fe.Message, nil)
		case fiber.StatusNotFound:
			return errorResp.NotFound(c, fe.Message)
		case fiber.StatusUnauthorized:
			return errorResp.Unauthorized(c, fe.Message)
		case fiber.StatusForbidden:
			return errorResp.Forbidden(c, fe.Message)
		case fiber.StatusConflict:
			return errorResp.Conflict(c, fe.Message)
		case fiber.StatusTooManyRequests:
			return errorResp.TooManyRequests(c, fe.Message)
		case fiber.StatusInternalServerError:
			return errorResp.InternalServerError(c, fe.Message)
		case fiber.StatusBadGateway:
			return errorResp.BadGateway(c, fe.Message)
		case fiber.StatusServiceUnavailable:
			return errorResp.ServiceUnavailable(c, fe.Message)
		case fiber.StatusGatewayTimeout:
			return errorResp.GatewayTimeout(c, fe.Message)
		default:
			log.Printf("Unhandled fiber error code: %d", fe.Code)
			return r.WriteError(c, fe.Code, "Error", "error", safeDetail(err), nil)
		}
	}
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorResp.NotFound(c, "Resource not found")
	}

	log.Println("Unhandled error fell through to fallback")
	return errorResp.InternalServerError(c, "An unexpected error occurred")
}

func safeDetail(err error) string {
	if isProd {
		log.Println("detail not shown in production:")
		return "" // ปิด detail ใน prod
	}
	return err.Error() // โชว์ใน dev
}

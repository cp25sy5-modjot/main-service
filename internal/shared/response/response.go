package response

import (
	"fmt"
	"os"
	"time"

	v "github.com/cp25sy5-modjot/main-service/internal/shared/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func WriteSuccess(c *fiber.Ctx, status int, data any, msg string) error {
	resp := Response{
		Method:    c.Method(),
		Path:      c.Path(),
		Status:    "success",
		Code:      status,
		Message:   msg,
		Data:      data,
		DraftID:   getDraftID(c),
		Timestamp: time.Now().UTC(),
	}
	LogSuccess(resp)
	return c.Status(status).JSON(resp)
}

func WriteError(c *fiber.Ctx, status int, msg, typ, detail string, fields []v.FieldError) error {
	resp := Response{
		Method:  c.Method(),
		Path:    c.Path(),
		Status:  "error",
		Code:    status,
		Message: msg,
		Error: &ErrorBody{
			Type:   typ,
			Detail: detail,
			Fields: fields,
		},
		DraftID:   getDraftID(c),
		Timestamp: time.Now().UTC(),
	}
	LogError(resp)
	return c.Status(status).JSON(resp)
}

func getDraftID(c *fiber.Ctx) string {
	if v := c.Locals("request_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func LogSuccess(resp Response) {
	logger.Info().
		Str("method", resp.Method).
		Str("path", resp.Path).
		Int("status_code", resp.Code).
		Str("detail", resp.Message).
		Str("draft_id", resp.DraftID).
		Msg("Success")
}

func LogError(resp Response) {
	logger.Error().
		Str("method", resp.Method).
		Str("path", resp.Path).
		Int("status_code", resp.Code).
		Str("error_message", resp.Message).
		Str("error_type", resp.Error.Type).
		Str("error_detail", resp.Error.Detail).
		Str("error_fields", fmt.Sprintf("%v", resp.Error.Fields)).
		Str("draft_id", resp.DraftID).
		Msg("Error")
}

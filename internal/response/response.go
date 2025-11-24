package response

import (
	"os"
	"time"

	v "github.com/cp25sy5-modjot/main-service/internal/validator"
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
		TraceID:   getTraceID(c),
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
		TraceID:   getTraceID(c),
		Timestamp: time.Now().UTC(),
	}
	LogError(resp)
	return c.Status(status).JSON(resp)
}

func getTraceID(c *fiber.Ctx) string {
	if v := c.Locals("request_id"); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func LogSuccess(resp Response) {
	logger.Info().
		Int("status_code", resp.Code).
		Str("message", resp.Message).
		Str("trace_id", resp.TraceID).
		Str("method", resp.Method).
		Str("path", resp.Path).
		Msg("HTTP Success Response")
}

func LogError(resp Response) {
	logger.Error().
		Int("status_code", resp.Code).
		Str("message", resp.Message).
		Str("error_type", resp.Error.Type).
		Str("error_detail", resp.Error.Detail).
		Str("trace_id", resp.TraceID).
		Str("method", resp.Method).
		Str("path", resp.Path).
		Msg("HTTP Error Response")
}

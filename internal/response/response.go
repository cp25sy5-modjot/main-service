package response

import (
	"log"
	"time"

	v "github.com/cp25sy5-modjot/main-service/internal/validator"
	"github.com/gofiber/fiber/v2"
)

func WriteSuccess(c *fiber.Ctx, status int, data any, msg string) error {
	resp := Response{
		Status:    "success",
		Code:      status,
		Message:   msg,
		Data:      data,
		TraceID:   getTraceID(c),
		Timestamp: time.Now().UTC(),
	}
	return c.Status(status).JSON(resp)
}

func WriteError(c *fiber.Ctx, status int, msg, typ, detail string, fields []v.FieldError) error {
	resp := Response{
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
	log.Println(resp)
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

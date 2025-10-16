package response

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status    string         `json:"status"`
	Code      int            `json:"code"`
	Message   string         `json:"message,omitempty"`
	Data      any            `json:"data,omitempty"`
	Error     *ErrorBody     `json:"error,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
	TraceID   string         `json:"trace_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type ErrorBody struct {
	Type   string       `json:"type,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Fields []FieldError `json:"fields,omitempty"`
}

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

func WriteError(c *fiber.Ctx, status int, msg, typ, detail string, fields []FieldError) error {
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

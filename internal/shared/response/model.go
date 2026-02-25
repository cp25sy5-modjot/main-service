package response

import (
	"time"

	v "github.com/cp25sy5-modjot/main-service/internal/shared/validator"
)

type Response struct {
	Method    string         `json:"method,omitempty"`
	Path      string         `json:"path,omitempty"`
	Status    string         `json:"status" valid:"required,oneof=success error"`
	Code      int            `json:"code" valid:"required"`
	Message   string         `json:"message,omitempty" valid:"required"`
	Data      any            `json:"data,omitempty"`
	Error     *ErrorBody     `json:"error,omitempty"`
	Meta      map[string]any `json:"meta,omitempty"`
	TraceID   string         `json:"trace_id,omitempty"`
	Timestamp time.Time      `json:"timestamp"`
}

type ErrorBody struct {
	Type             string         `json:"type,omitempty"`
	Detail           string         `json:"detail,omitempty"`
	Fields           []v.FieldError `json:"fields,omitempty"`
	RemainingSeconds int            `json:"remaining_seconds,omitempty"`
}

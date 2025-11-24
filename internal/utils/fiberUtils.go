package utils

import (
	v "github.com/cp25sy5-modjot/main-service/internal/validator"
	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, req interface{}) error {
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
	}
	return nil
}

func ValidateStruct(c *fiber.Ctx, req interface{}) error {
	if err := v.Validator().Struct(req); err != nil {
		return UnprocessableEntity("Validation Failed", err)
	}
	return nil
}

func ParseBodyAndValidate(c *fiber.Ctx, req interface{}) error {
	if err := ParseBody(c, req); err != nil {
		return err
	}
	if err := ValidateStruct(c, req); err != nil {
		return err
	}
	return nil
}

func UnprocessableEntity(detail string, err error) error {
	return &ValidationError{
		OriginalErr: err,
		Message:     detail,
	}
}

// Note: You might want a custom error type for validation to pass the fields
// Here's an example:
type ValidationError struct {
	OriginalErr error
	Message     string
}

func (e *ValidationError) Error() string {
	return e.Message
}

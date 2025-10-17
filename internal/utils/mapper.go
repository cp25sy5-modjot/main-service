package utils

import (
	"dario.cat/mergo"
	errResp "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func MapStructs(c *fiber.Ctx, src interface{}, dest interface{}) error {
	if err := copier.Copy(dest, src); err != nil {
		return errResp.InternalServerError(c, "Failed to map structures")
	}
	return nil
}

func MapNonNilStructs(c *fiber.Ctx, src interface{}, dest interface{}) error {
	if err := mergo.Merge(dest, src); err != nil {
		return errResp.InternalServerError(c, "Failed to map non-nil structures")
	}
	return nil
}

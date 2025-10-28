package utils

import (
	"dario.cat/mergo"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func MapStructs(src interface{}, dest interface{}) error {
	if err := copier.Copy(dest, src); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to map structures")
	}
	return nil
}

func MapNonNilStructs(src interface{}, dest interface{}) error {
	if err := mergo.Merge(dest, src); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to map non-nil structures")
	}
	return nil
}

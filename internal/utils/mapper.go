package utils

import (
	"dario.cat/mergo"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

//for map response structs (shallow copy)
func MapStructs(src, dest interface{}) error {
	if err := copier.Copy(dest, src); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to map structures")
	}
	return nil
}

// merge structs, dest gets overwritten by src fields that exist in src (work only same struct types)
func MergeStructs(src, dest interface{}) error {
	// dest gets overwritten by src fields that exist in src
	if err := mergo.Merge(dest, src, mergo.WithOverride); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to merge structures")
	}
	return nil
}

// map slice of structs with deeper copy
func MapSliceOfStructs(src, dest interface{}) error {
	if err := copier.CopyWithOption(dest, src, copier.Option{DeepCopy: true}); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to map slice of structures")
	}
	return nil
}

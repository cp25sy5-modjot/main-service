package success

import (
	"github.com/cp25sy5-modjot/main-service/internal/response"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

// OK → 200
func OK(c *fiber.Ctx, data any, msg ...string) error {
	m := utils.FirstOrEmpty(msg, "OK")
	return response.WriteSuccess(c, fiber.StatusOK, data, m)
}

// Created → 201
func Created(c *fiber.Ctx, data any, msg ...string) error {
	m := utils.FirstOrEmpty(msg, "Created successfully")
	return response.WriteSuccess(c, fiber.StatusCreated, data, m)
}

// Accepted → 202
func Accepted(c *fiber.Ctx, data any, msg ...string) error {
	m := utils.FirstOrEmpty(msg, "Accepted")
	return response.WriteSuccess(c, fiber.StatusAccepted, data, m)
}

// NoContent → 204
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// OKWithMeta → 200 + meta (เช่น pagination)
// func OKWithMeta(c *fiber.Ctx, data any, meta map[string]any, msg ...string) error {
// 	m := firstOrEmpty(msg, "OK")
// 	env := response.Response{
// 		Status:    "success",
// 		Code:      fiber.StatusOK,
// 		Message:   m,
// 		Data:      data,
// 		Meta:      meta,
// 		TraceID:   response.getTraceID(c),
// 		Timestamp: time.Now().UTC(),
// 	}
// 	return c.Status(fiber.StatusOK).JSON(env)
// }

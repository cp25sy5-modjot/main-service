package draft

import (
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	"github.com/gofiber/fiber/v2"
		mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"

)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetDraft(c *fiber.Ctx) error {

	traceID := c.Params("traceID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	d, err := h.service.GetDraft(c.Context(), traceID, userID)
	if err != nil {
		return fiber.NewError(404, "draft not found")
	}

	if d.UserID != userID {
		return fiber.NewError(403, "not owner")
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   d,
	})
}

func (h *Handler) ListDraft(c *fiber.Ctx) error {

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	list, err := h.service.ListDraft(c.Context(), userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   list,
	})
}

func (h *Handler) Update(c *fiber.Ctx) error {

	traceID := c.Params("traceID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req ConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	d, err := h.service.UpdateDraft(c.Context(), traceID, userID, req)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   d,
	})
}

func (h *Handler) Confirm(c *fiber.Ctx) error {

	traceID := c.Params("traceID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req ConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	tx, err := h.service.ConfirmDraft(
		c.Context(),
		traceID,
		userID,
		req,
	)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.JSON(fiber.Map{
		"status":      "success",
		"transaction": mapper.BuildTransactionResponse(tx),
	})
}

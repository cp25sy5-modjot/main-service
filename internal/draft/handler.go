package draft

import (
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetDraft(c *fiber.Ctx) error {

	draftID := c.Params("draftID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	d, err := h.service.GetDraftWithCategory(c.Context(), draftID, userID)
	if err != nil {
		return fiber.NewError(404, "draft not found")
	}

	return sresp.OK(c, d, "draft retrieved successfully")
}

func (h *Handler) ListDraft(c *fiber.Ctx) error {

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	list, err := h.service.ListDraftWithCategory(c.Context(), userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return sresp.OK(c, list, "drafts retrieved successfully")
}

func (h *Handler) Update(c *fiber.Ctx) error {

	draftID := c.Params("draftID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req ConfirmRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid body")
	}

	d, err := h.service.UpdateDraft(c.Context(), draftID, userID, req)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return sresp.OK(c, d, "draft updated successfully")
}

func (h *Handler) Confirm(c *fiber.Ctx) error {
	draftID := c.Params("draftID")
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
		draftID,
		userID,
		req,
	)

	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return sresp.OK(c, mapper.BuildTransactionResponse(tx), "draft confirmed successfully")
}

func (h *Handler) GetDraftStats(c *fiber.Ctx) error {

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	stats, err := h.service.GetDraftStats(c.Context(), userID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}

	return sresp.OK(c, stats, "draft stats retrieved successfully")
}

func (h *Handler) GetDraftImageURL(c *fiber.Ctx) error {

	draftID := c.Params("draftID")
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	url, err := h.service.GetDraftImageURL(
		c.Context(),
		draftID,
		userID,
	)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	return c.JSON(fiber.Map{
		"url": url,
	})
}

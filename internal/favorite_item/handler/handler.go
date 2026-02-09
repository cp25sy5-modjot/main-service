package favhandler

import (
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	svc "github.com/cp25sy5-modjot/main-service/internal/favorite_item/service"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	mapper "github.com/cp25sy5-modjot/main-service/internal/mapper"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service svc.Service
}

func NewHandler(service svc.Service) *Handler {
	return &Handler{service}
}

// POST /favorites
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req m.FavoriteItemInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	var input = mapper.ParseFavoriteItemInsertReqToServiceInput(userID, &req)

	fav, err := h.service.Create(input)
	if err != nil {
		return err
	}

	return sresp.Created(c, mapper.BuildFavoriteItemResponse(fav), "Favorite item created successfully")
}

// PUT /favorites
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	favID := c.Params("id")
	if favID == "" {
		return fiber.ErrBadRequest
	}

	var req m.FavoriteItemUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	input := mapper.ParseFavoriteItemUpdateReqToServiceInput(userID, favID, &req)

	updatedFav, err := h.service.Update(input)
	if err != nil {
		return err
	}

	return sresp.OK(
		c,
		mapper.BuildFavoriteItemResponse(updatedFav),
		"Favorite item updated successfully",
	)
}

// DELETE /favorites/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	favID := c.Params("id")
	if favID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "favorite_id parameter is required")
	}

	if err := h.service.Delete(userID, favID); err != nil {
		return err
	}

	return sresp.NoContent(c)
}

// GET /favorites
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	favs, err := h.service.GetAll(userID)
	if err != nil {
		return err
	}

	var favResList []*m.FavoriteItemRes
	for _, fav := range favs {
		favResList = append(favResList, mapper.BuildFavoriteItemResponse(fav))
	}

	return sresp.OK(c, favResList, "Favorite items retrieved successfully")
}

// GET /favorites/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	favID := c.Params("id")
	if favID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "favorite_id parameter is required")
	}

	fav, err := h.service.GetByID(userID, favID)
	if err != nil {
		return err
	}

	return sresp.OK(
		c,
		mapper.BuildFavoriteItemResponse(fav),
		"Favorite item retrieved successfully",
	)
}

// PUT /favorites/reorder
func (h *Handler) ReOrder(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req m.FavoriteItemReOrderReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	input := mapper.ParseFavoriteItemReOrderReqToServiceInput(userID, &req)

	if err := h.service.ReOrder(input); err != nil {
		return err
	}

	return sresp.NoContent(c)
}

package category

import (
	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	successResp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *catSvc.Service
}

func NewHandler(service *catSvc.Service) *Handler {
	return &Handler{service}
}

// POST /category
func (h *Handler) Create(c *fiber.Ctx) error {
	var req m.CategoryReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	cate := &e.Category{
		CategoryName: req.CategoryName,
		Budget:       req.Budget,
		UserID:       userID,
	}

	createdCate, err := h.service.Create(cate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var cateRes m.CategoryRes
	utils.MapStructs(createdCate, &cateRes)
	return successResp.Created(c, cateRes, "Category created successfully")
}

// GET /categories
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	categories, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var cateRes []m.CategoryRes
	utils.MapStructs(&categories, &cateRes)
	return successResp.OK(c, cateRes, "Categories retrieved successfully")
}

// PUT /category/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req m.CategoryUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	params := &m.CategorySearchParams{
		CategoryID: c.Params("id"),
		UserID:     userID,
	}

	updatedCategory, err := h.service.Update(params, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var cateRes m.CategoryRes
	utils.MapStructs(updatedCategory, &cateRes)
	return successResp.OK(c, cateRes, "Category updated successfully")
}

// DELETE /category/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	params := &m.CategorySearchParams{
		CategoryID: c.Params("id"),
		UserID:     userID,
	}

	if err := h.service.Delete(params); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return successResp.OK(c, nil, "Category deleted successfully")
}

package category

import (
	"strconv"

	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
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

	isIncludeTransactions, err := isIncludeTransactions(c)
	if err != nil {
		return err
	}

	if isIncludeTransactions {
		categories, err := h.service.GetAllByUserIDWithTransactions(userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve categories with transactions")
		}
		return successResp.OK(c, categories, "Categories with transactions retrieved successfully")
	}

	categories, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve categories")
	}
	return successResp.OK(c, categories, "Categories retrieved successfully")
}

// GET /category/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	params := &m.CategorySearchParams{
		CategoryID: c.Params("id"),
		UserID:     userID,
	}
	
	isIncludeTransactions, err := isIncludeTransactions(c)
	if err != nil {
		return err
	}

	if isIncludeTransactions {
		category, err := h.service.GetByIDWithTransactions(params)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return successResp.OK(c, category, "Category with transactions retrieved successfully")
	}

	category, err := h.service.GetByID(params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return successResp.OK(c, category, "Category retrieved successfully")
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

	category, err := h.service.Update(params, &req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return successResp.OK(c, category, "Category updated successfully")
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

// utils
func isIncludeTransactions(c *fiber.Ctx) (bool, error) {
	includeTransactionsStr := c.Query("includeTransactions", "false")
	includeTransactions, err := strconv.ParseBool(includeTransactionsStr)
	if err != nil {
		return false, fiber.NewError(fiber.StatusBadRequest, "includeTransactions must be a boolean (true/false)")
	}
	return includeTransactions, nil
}

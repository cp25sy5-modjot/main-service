package categoryhandler

import (
	"strconv"

	categorysvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service categorysvc.Service
}

func NewHandler(service categorysvc.Service) *Handler {
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

	cate := &categorysvc.CategoryCreateInput{
		CategoryName: req.CategoryName,
		Budget:       req.Budget,
		ColorCode:    req.ColorCode,
	}

	createdCate, err := h.service.Create(userID, cate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return sresp.Created(c, buildCategoryResponse(createdCate), "Category created successfully")
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
		date := c.Query("date")
		filter := &m.TransactionFilter{
			Date: utils.ConvertStringToTime(date),
		}
		categories, err := h.service.GetAllByUserIDWithTransactions(userID, filter)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve categories with transactions")
		}
		return sresp.OK(c, categories, "Categories with transactions retrieved successfully")
	}

	categories, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to retrieve categories")
	}
	return sresp.OK(c, buildCategoryResponses(categories), "Categories retrieved successfully")
}

// GET /category/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	categoryID := c.Params("id")
	params := &m.CategorySearchParams{
		CategoryID: &categoryID,
		UserID:     userID,
	}

	isIncludeTransactions, err := isIncludeTransactions(c)
	if err != nil {
		return err
	}

	if isIncludeTransactions {
		date := c.Query("date")
		filter := &m.TransactionFilter{
			Date: utils.ConvertStringToTime(date),
		}
		category, err := h.service.GetByIDWithTransactions(params, filter)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return sresp.OK(c, buildCategoryResponse(category), "Category with transactions retrieved successfully")
	}

	category, err := h.service.GetByID(params)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sresp.OK(c, buildCategoryResponse(category), "Category retrieved successfully")
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

	categoryID := c.Params("id")
	params := &m.CategorySearchParams{
		CategoryID: &categoryID,
		UserID:     userID,
	}

	update := &categorysvc.CategoryUpdateInput{
		CategoryName: req.CategoryName,
		Budget:       req.Budget,
		ColorCode:    req.ColorCode,
	}

	category, err := h.service.Update(params, update)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return sresp.OK(c, buildCategoryResponse(category), "Category updated successfully")
}

// DELETE /category/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	categoryID := c.Params("id")
	params := &m.CategorySearchParams{
		CategoryID: &categoryID,
		UserID:     userID,
	}

	if err := h.service.Delete(params); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sresp.OK(c, nil, "Category deleted successfully")
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

func buildCategoryResponse(cat *e.Category) *m.CategoryRes {
	return &m.CategoryRes{
		CategoryID:   &cat.CategoryID,
		CategoryName: cat.CategoryName,
		Budget:       cat.Budget,
		ColorCode:    cat.ColorCode,
		CreatedAt:    cat.CreatedAt,
	}
}

func buildCategoryResponses(categories []e.Category) []m.CategoryRes {
	categoryResponses := make([]m.CategoryRes, 0, len(categories))
	for _, cat := range categories {
		res := buildCategoryResponse(&cat)
		categoryResponses = append(categoryResponses, *res)
	}
	return categoryResponses
}

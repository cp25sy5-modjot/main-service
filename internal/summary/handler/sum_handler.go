package summaryhandler

import (
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/cp25sy5-modjot/main-service/internal/summary/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service}
}

// GET /summary?period=week|month|year
func (h *Handler) GetExpenseSummary(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	period := c.Query("period")
	if period == "" {
		return fiber.ErrBadRequest
	}

	summary, err := h.service.GetExpenseSummary(c.Context(), userID, service.Period(period))
	if err != nil {
		return err
	}

	return sresp.OK(c, summary, "Expense summary retrieved successfully")
}

// GET /summary/category?period=week|month|year&date=2024-01-01
func (h *Handler) GetCategorySummary(c *fiber.Ctx) error {

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var q m.CategorySummaryQuery

	if err := utils.ParseQueryAndValidate(c, &q); err != nil {
		return err
	}

	period := service.Period(q.Period)
	date := utils.ConvertStringToTime(q.Date)

	summary, err := h.service.GetCategorySummary(
		c.Context(),
		userID,
		period,
		date,
	)
	if err != nil {
		return err
	}

	return sresp.OK(c, summary, "Category summary retrieved successfully")
}

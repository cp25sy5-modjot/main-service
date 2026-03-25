package summaryhandler_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	summaryhandler "github.com/cp25sy5-modjot/main-service/internal/summary/handler"
	"github.com/cp25sy5-modjot/main-service/internal/summary/mocks"
	service "github.com/cp25sy5-modjot/main-service/internal/summary/service"
)

func TestSummaryHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	// ===================== EXPENSE SUMMARY =====================
	t.Run("Expense Summary Success", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				if e, ok := err.(*fiber.Error); ok {
					return c.Status(e.Code).JSON(fiber.Map{"error": e.Message})
				}
				return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
			},
		})

		mockSvc := mocks.NewMockService(t)
		h := summaryhandler.NewHandler(mockSvc)

		app.Get("/summary", h.GetExpenseSummary)

		mockSvc.EXPECT().
			GetExpenseSummary(mock.Anything, "u-1", service.Period("month")).
			Return(m.ExpenseSummaryRes{}, nil)

		req := httptest.NewRequest("GET", "/summary?period=month", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Expense Summary Missing Period", func(t *testing.T) {
		app := fiber.New()

		mockSvc := mocks.NewMockService(t)
		h := summaryhandler.NewHandler(mockSvc)

		app.Get("/summary", h.GetExpenseSummary)

		req := httptest.NewRequest("GET", "/summary", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("Expense Summary Error", func(t *testing.T) {
		app := fiber.New()

		mockSvc := mocks.NewMockService(t)
		h := summaryhandler.NewHandler(mockSvc)

		app.Get("/summary", h.GetExpenseSummary)

		mockSvc.EXPECT().
			GetExpenseSummary(mock.Anything, "u-1", service.Period("month")).
			Return(m.ExpenseSummaryRes{}, errors.New("error"))

		req := httptest.NewRequest("GET", "/summary?period=month", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== CATEGORY SUMMARY =====================
	t.Run("Category Summary Success", func(t *testing.T) {
		app := fiber.New()

		mockSvc := mocks.NewMockService(t)
		h := summaryhandler.NewHandler(mockSvc)

		app.Get("/summary/category", h.GetCategorySummary)

		mockSvc.EXPECT().
			GetCategorySummary(
				mock.Anything,
				"u-1",
				service.Period("month"),
				mock.Anything,
			).
			Return(m.CategorySummaryRes{}, nil)

		req := httptest.NewRequest("GET", "/summary/category?period=month&date=2024-01-01", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Category Summary Error", func(t *testing.T) {
		app := fiber.New()

		mockSvc := mocks.NewMockService(t)
		h := summaryhandler.NewHandler(mockSvc)

		app.Get("/summary/category", h.GetCategorySummary)

		mockSvc.EXPECT().
			GetCategorySummary(
				mock.Anything,
				"u-1",
				service.Period("month"),
				mock.Anything,
			).
			Return(m.CategorySummaryRes{}, errors.New("error"))

		req := httptest.NewRequest("GET", "/summary/category?period=month&date=2024-01-01", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
package overviewhandler_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	overviewhandler "github.com/cp25sy5-modjot/main-service/internal/overview/handler"
	"github.com/cp25sy5-modjot/main-service/internal/overview/mocks"
)

func TestOverviewHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	// ===================== SUCCESS WITH DATE =====================
	t.Run("Success With Date", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := overviewhandler.NewHandler(mockSvc)
		app.Get("/overview", h.GetOverview)

		mockSvc.EXPECT().
			GetOverview("u-1", mock.AnythingOfType("time.Time")).
			Return(&m.OverviewResponse{}, nil)

		req := httptest.NewRequest("GET", "/overview?date=2024-01-01", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== SUCCESS NO DATE =====================
	t.Run("Success No Date", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := overviewhandler.NewHandler(mockSvc)
		app.Get("/overview", h.GetOverview)

		mockSvc.EXPECT().
			GetOverview("u-1", mock.AnythingOfType("time.Time")).
			Return(&m.OverviewResponse{}, nil)

		req := httptest.NewRequest("GET", "/overview", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== INVALID DATE =====================
	t.Run("Invalid Date", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := overviewhandler.NewHandler(mockSvc)
		app.Get("/overview", h.GetOverview)

		req := httptest.NewRequest("GET", "/overview?date=2024/01/01", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	// ===================== SERVICE ERROR =====================
	t.Run("Service Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := overviewhandler.NewHandler(mockSvc)
		app.Get("/overview", h.GetOverview)

		mockSvc.EXPECT().
			GetOverview("u-1", mock.AnythingOfType("time.Time")).
			Return(nil, errors.New("service error"))

		req := httptest.NewRequest("GET", "/overview?date=2024-01-01", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
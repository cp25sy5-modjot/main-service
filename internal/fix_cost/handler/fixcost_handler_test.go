package fixcosthandler_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	fixcosthandler "github.com/cp25sy5-modjot/main-service/internal/fix_cost/handler"
	"github.com/cp25sy5-modjot/main-service/internal/fix_cost/mocks"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

func TestFixCostHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	validCreateBody := `{
		"title": "Netflix",
		"price": 199,
		"category_id": "c-1",
		"interval_type": "monthly",
		"interval_value": 1,
		"start_date": "2024-01-01T00:00:00Z"
	}`

	validUpdateBody := `{
		"title": "Spotify",
		"price": 100,
		"category_id": "c-1",
		"interval_type": "monthly",
		"interval_value": 1,
		"start_date": "2024-01-01T00:00:00Z",
		"status": "active"
	}`

	// ===================== CREATE =====================
	t.Run("Create Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Post("/fix_cost", h.Create)

		mockSvc.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(&e.FixCost{Title: "Netflix"}, nil)

		req := httptest.NewRequest("POST", "/fix_cost", bytes.NewBufferString(validCreateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 201, resp.StatusCode)
	})

	t.Run("Create Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Post("/fix_cost", h.Create)

		mockSvc.EXPECT().
			Create(mock.Anything, mock.Anything).
			Return(nil, errors.New("error"))

		req := httptest.NewRequest("POST", "/fix_cost", bytes.NewBufferString(validCreateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== UPDATE =====================
	t.Run("Update Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Put("/fix_cost/:id", h.Update)

		mockSvc.EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(&e.FixCost{Title: "Spotify"}, nil)

		req := httptest.NewRequest("PUT", "/fix_cost/fc-1", bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Update BadRequest", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Put("/fix_cost/:id", h.Update)

		req := httptest.NewRequest("PUT", "/fix_cost/%20", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	// ===================== DELETE =====================
	t.Run("Delete Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Delete("/fix_cost/:id", h.Delete)

		mockSvc.EXPECT().
			Delete(mock.Anything, "fc-1", "u-1").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/fix_cost/fc-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 204, resp.StatusCode)
	})

	t.Run("Delete Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Delete("/fix_cost/:id", h.Delete)

		mockSvc.EXPECT().
			Delete(mock.Anything, "fc-1", "u-1").
			Return(errors.New("error"))

		req := httptest.NewRequest("DELETE", "/fix_cost/fc-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== GET BY ID =====================
	t.Run("GetByID Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Get("/fix_cost/:id", h.GetByID)

		mockSvc.EXPECT().
			GetByID(mock.Anything, "fc-1", "u-1").
			Return(&e.FixCost{Title: "Netflix"}, nil)

		req := httptest.NewRequest("GET", "/fix_cost/fc-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetByID Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Get("/fix_cost/:id", h.GetByID)

		mockSvc.EXPECT().
			GetByID(mock.Anything, "fc-1", "u-1").
			Return(nil, errors.New("error"))

		req := httptest.NewRequest("GET", "/fix_cost/fc-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== GET ALL =====================
	t.Run("GetAll Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Get("/fix_cost", h.GetAllByUserID)

		mockSvc.EXPECT().
			GetAllByUserID(mock.Anything, "u-1").
			Return([]*e.FixCost{
				{Title: "Netflix"},
				{Title: "Spotify"},
			}, nil)

		req := httptest.NewRequest("GET", "/fix_cost", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetAll Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := fixcosthandler.NewHandler(mockSvc)
		app.Get("/fix_cost", h.GetAllByUserID)

		mockSvc.EXPECT().
			GetAllByUserID(mock.Anything, "u-1").
			Return(nil, errors.New("error"))

		req := httptest.NewRequest("GET", "/fix_cost", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
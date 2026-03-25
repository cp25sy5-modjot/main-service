package categoryhandler_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	categoryhandler "github.com/cp25sy5-modjot/main-service/internal/category/handler"
	"github.com/cp25sy5-modjot/main-service/internal/category/mocks"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

func TestCategoryHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	validCreateBody := `{
		"category_name": "Food",
		"budget": 1000,
		"color_code": "#FF0000",
		"icon": "🍔"
	}`

	validUpdateBody := `{
		"category_name": "Food Updated",
		"budget": 2000,
		"color_code": "#00FF00",
		"icon": "🍕"
	}`

	// ===================== CREATE =====================
	t.Run("Create Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Post("/category", h.Create)

		mockSvc.EXPECT().
			Create("u-1", mock.Anything).
			Return(&e.Category{CategoryName: "Food"}, nil)

		req := httptest.NewRequest("POST", "/category", bytes.NewBufferString(validCreateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 201, resp.StatusCode)
	})

	// ===================== GET ALL =====================
	t.Run("GetAll Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Get("/categories", h.GetAll)

		mockSvc.EXPECT().
			GetAllByUserID("u-1").
			Return([]e.Category{{CategoryName: "Food"}}, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetAll With Transactions", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Get("/categories", h.GetAll)

		mockSvc.EXPECT().
			GetAllByUserIDWithTransactions("u-1", mock.Anything).
			Return([]e.Category{{CategoryName: "Food"}}, nil)

		req := httptest.NewRequest("GET", "/categories?includeTransactions=true", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetAll Invalid Query", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Get("/categories", h.GetAll)

		req := httptest.NewRequest("GET", "/categories?includeTransactions=abc", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	// ===================== GET BY ID =====================
	t.Run("GetByID Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Get("/category/:id", h.GetByID)

		mockSvc.EXPECT().
			GetByID(mock.Anything).
			Return(&e.Category{CategoryName: "Food"}, nil)

		req := httptest.NewRequest("GET", "/category/c-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetByID With Transactions", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Get("/category/:id", h.GetByID)

		mockSvc.EXPECT().
			GetByIDWithTransactions(mock.Anything, mock.Anything).
			Return(&e.Category{CategoryName: "Food"}, nil)

		req := httptest.NewRequest("GET", "/category/c-1?includeTransactions=true", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== UPDATE =====================
	t.Run("Update Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Put("/category/:id", h.Update)

		mockSvc.EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(&e.Category{CategoryName: "Updated"}, nil)

		req := httptest.NewRequest("PUT", "/category/c-1", bytes.NewBufferString(validUpdateBody))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== DELETE =====================
	t.Run("Delete Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := categoryhandler.NewHandler(mockSvc)
		app.Delete("/category/:id", h.Delete)

		mockSvc.EXPECT().
			Delete(mock.Anything).
			Return(nil)

		req := httptest.NewRequest("DELETE", "/category/c-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})
}
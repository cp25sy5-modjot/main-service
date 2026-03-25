package favhandler_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	favhandler "github.com/cp25sy5-modjot/main-service/internal/favorite_item/handler"
	"github.com/cp25sy5-modjot/main-service/internal/favorite_item/mocks"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

func TestFavoriteHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	validUpdateBody := `{
		"title": "Tea",
		"price": 30,
		"category_id": "550e8400-e29b-41d4-a716-446655440000"
	}`

	validReorderBody := `{
		"reorder_list": [
			{
				"favorite_id": "550e8400-e29b-41d4-a716-446655440000",
				"position": 1
			},
			{
				"favorite_id": "550e8400-e29b-41d4-a716-446655440001",
				"position": 2
			}
		]
	}`

	// ===================== UPDATE =====================
	t.Run("Update Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Put("/favorites/:id", h.Update)

		mockSvc.EXPECT().
			Update(mock.Anything).
			Return(&e.FavoriteItem{Title: "Tea"}, nil)

		req := httptest.NewRequest(
			"PUT",
			"/favorites/fav-1",
			bytes.NewBufferString(validUpdateBody),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Update BadRequest", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Put("/favorites/:id", h.Update)

		req := httptest.NewRequest("PUT", "/favorites/%20", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})

	// ===================== DELETE =====================
	t.Run("Delete Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Delete("/favorites/:id", h.Delete)

		mockSvc.EXPECT().
			Delete("u-1", "fav-1").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/favorites/fav-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 204, resp.StatusCode)
	})

	// ===================== GET ALL =====================
	t.Run("GetAll Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Get("/favorites", h.GetAll)

		mockSvc.EXPECT().
			GetAll("u-1").
			Return([]*e.FavoriteItem{
				{Title: "Coffee"},
			}, nil)

		req := httptest.NewRequest("GET", "/favorites", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== GET BY ID =====================
	t.Run("GetByID Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Get("/favorites/:id", h.GetByID)

		mockSvc.EXPECT().
			GetByID("u-1", "fav-1").
			Return(&e.FavoriteItem{Title: "Coffee"}, nil)

		req := httptest.NewRequest("GET", "/favorites/fav-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== REORDER =====================
	t.Run("ReOrder Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := favhandler.NewHandler(mockSvc)
		app.Put("/favorites/reorder", h.ReOrder)

		mockSvc.EXPECT().
			ReOrder(mock.Anything).
			Return(nil)

		req := httptest.NewRequest(
			"PUT",
			"/favorites/reorder",
			bytes.NewBufferString(validReorderBody),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 204, resp.StatusCode)
	})
}
package userhandler_test

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"
	userhandler "github.com/cp25sy5-modjot/main-service/internal/user/handler"
	"github.com/cp25sy5-modjot/main-service/internal/user/mocks"
)

func TestUserHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	// ===================== GET SELF =====================
	t.Run("GetSelf Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Get("/user", h.GetSelf)

		mockSvc.EXPECT().
			GetByID("u-1").
			Return(&e.User{Name: "James"}, nil)

		req := httptest.NewRequest("GET", "/user", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetSelf Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Get("/user", h.GetSelf)

		mockSvc.EXPECT().
			GetByID("u-1").
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest("GET", "/user", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== CREATE =====================
	t.Run("Create Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Post("/user", h.Create)

		body := `{
			"name": "James",
			"userBinding": {
				"googleID": "g-123"
			}
		}`

		mockSvc.EXPECT().
			Create(mock.Anything).
			Return(&e.User{Name: "James"}, nil)

		req := httptest.NewRequest(
			"POST",
			"/user",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 201, resp.StatusCode)
	})

	// ===================== UPDATE =====================
	t.Run("Update Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Put("/user", h.Update)

		body := `{
			"name": "NewName"
		}`

		mockSvc.EXPECT().
			Update("u-1", mock.Anything).
			Return(&e.User{Name: "NewName"}, nil)

		req := httptest.NewRequest(
			"PUT",
			"/user",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	// ===================== DELETE =====================
	t.Run("Delete Soft", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Delete("/user", h.Delete)

		mockSvc.EXPECT().
			SoftDelete("u-1").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/user?mode=soft", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Delete Hard", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Delete("/user", h.Delete)

		mockSvc.EXPECT().
			Delete("u-1").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/user?mode=hard", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Delete Invalid Mode", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := userhandler.NewHandler(mockSvc)
		app.Delete("/user", h.Delete)

		req := httptest.NewRequest("DELETE", "/user?mode=unknown", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	})
}
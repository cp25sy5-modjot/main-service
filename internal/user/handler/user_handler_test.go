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
	"github.com/cp25sy5-modjot/main-service/internal/user/mocks")

//
// 🔥 MOCK JWT (override behavior)
//

// create a fake function variable (you need this in your jwt package)
// see note below 👇
var fakeGetUserID = func(c *fiber.Ctx) (string, error) {
	return "u-1", nil
}

//
// 🧪 TEST: GetSelf
//

func TestGetSelf_Success(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	// 👇 override jwt behavior
	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Get("/user", h.GetSelf)

	mockSvc.On("GetByID", "u-1").
		Return(&e.User{Name: "James"}, nil)

	req := httptest.NewRequest("GET", "/user", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockSvc.AssertExpectations(t)
}

//
// 🧪 TEST: GetSelf error
//

func TestGetSelf_Error(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Get("/user", h.GetSelf)

	mockSvc.On("GetByID", "u-1").
		Return(nil, errors.New("db error"))

	req := httptest.NewRequest("GET", "/user", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

//
// 🧪 TEST: Create (no jwt needed)
//

func TestCreate_Success(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	app.Post("/user", h.Create)

	body := `{
		"name": "James",
		"userBinding": {
			"googleID": "g-123"
		}
	}`

	mockSvc.On("Create", mock.Anything).
		Return(&e.User{Name: "James"}, nil)

	req := httptest.NewRequest("POST", "/user", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
}

//
// 🧪 TEST: Update
//

func TestUpdate_Success(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Put("/user", h.Update)

	body := `{
		"name": "NewName"
	}`

	mockSvc.On("Update", "u-1", mock.Anything).
		Return(&e.User{Name: "NewName"}, nil)

	req := httptest.NewRequest("PUT", "/user", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// 🧪 TEST: Delete soft
//

func TestDelete_Soft(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Delete("/user", h.Delete)

	mockSvc.On("SoftDelete", "u-1").Return(nil)

	req := httptest.NewRequest("DELETE", "/user?mode=soft", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// 🧪 TEST: Delete hard
//

func TestDelete_Hard(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Delete("/user", h.Delete)

	mockSvc.On("Delete", "u-1").Return(nil)

	req := httptest.NewRequest("DELETE", "/user?mode=hard", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
}

//
// 🧪 TEST: Delete invalid mode
//

func TestDelete_InvalidMode(t *testing.T) {
	app := fiber.New()

	mockSvc := new(mocks.Service)
	h := userhandler.NewHandler(mockSvc)

	jwt.GetUserIDFromClaims = fakeGetUserID

	app.Delete("/user", h.Delete)

	req := httptest.NewRequest("DELETE", "/user?mode=unknown", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}

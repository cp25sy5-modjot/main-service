package drafthandler_test

import (
	"bytes"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	drafthandler "github.com/cp25sy5-modjot/main-service/internal/draft/handler"
	"github.com/cp25sy5-modjot/main-service/internal/draft/mocks"
	jwt "github.com/cp25sy5-modjot/main-service/internal/jwt"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	model "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func TestHandler_All(t *testing.T) {

	// fake jwt
	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	// ===================== GetDraft SUCCESS =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Get("/draft/:draftID", h.GetDraft)

		mockSvc.EXPECT().
			GetDraftWithCategory(mock.Anything, "d-1", "u-1").
			Return(&model.DraftRes{}, nil)

		req := httptest.NewRequest("GET", "/draft/d-1", nil)
		resp, _ := app.Test(req)

		body, _ := io.ReadAll(resp.Body)

		assert.Equal(t, 200, resp.StatusCode)
		assert.NotEmpty(t, body)
	}

	// ===================== GetDraft NOT FOUND =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Get("/draft/:draftID", h.GetDraft)

		mockSvc.EXPECT().
			GetDraftWithCategory(mock.Anything, "d-1", "u-1").
			Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/draft/d-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 404, resp.StatusCode)
	}

	// ===================== LIST =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Get("/draft", h.ListDraft)

		mockSvc.EXPECT().
			ListDraftWithCategory(mock.Anything, "u-1").
			Return([]model.DraftRes{}, nil)

		req := httptest.NewRequest("GET", "/draft", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}

	// ===================== UPDATE SUCCESS =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Put("/draft/:draftID", h.Update)

		mockSvc.EXPECT().
			UpdateDraft(
				mock.Anything,
				"d-1",
				"u-1",
				mock.Anything,
			).
			Return(&model.DraftTxn{}, nil)

		req := httptest.NewRequest(
			"PUT",
			"/draft/d-1",
			bytes.NewBufferString(`{"items":[{"price":100}]}`),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}

	// ===================== UPDATE BAD BODY =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Put("/draft/:draftID", h.Update)

		req := httptest.NewRequest("PUT", "/draft/d-1", bytes.NewBufferString("bad"))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	}

	// ===================== CONFIRM SUCCESS =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Post("/draft/:draftID/confirm", h.Confirm)

		mockSvc.EXPECT().
			ConfirmDraft(mock.Anything, "d-1", "u-1", mock.Anything).
			Return(&e.Transaction{}, nil)

		req := httptest.NewRequest(
			"POST",
			"/draft/d-1/confirm",
			bytes.NewBufferString(`{"items":[{"price":100}]}`),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}

	// ===================== CONFIRM ERROR =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Post("/draft/:draftID/confirm", h.Confirm)

		mockSvc.EXPECT().
			ConfirmDraft(mock.Anything, "d-1", "u-1", mock.Anything).
			Return(nil, errors.New("fail"))

		req := httptest.NewRequest(
			"POST",
			"/draft/d-1/confirm",
			bytes.NewBufferString(`{"items":[{"price":100}]}`),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 400, resp.StatusCode)
	}

	// ===================== DELETE =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Delete("/draft/:draftID", h.DeleteDraft)

		mockSvc.EXPECT().
			DeleteDraft(mock.Anything, "d-1", "u-1").
			Return(nil)

		req := httptest.NewRequest("DELETE", "/draft/d-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}

	// ===================== STATS =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Get("/draft/stats", h.GetDraftStats)

		mockSvc.EXPECT().
			GetDraftStats(mock.Anything, "u-1").
			Return(&model.DraftStats{}, nil)

		req := httptest.NewRequest("GET", "/draft/stats", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}

	// ===================== IMAGE =====================
	{
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := drafthandler.NewHandler(mockSvc)
		app.Get("/draft/:draftID/image", h.GetDraftImageURL)

		mockSvc.EXPECT().
			GetDraftImageURL(mock.Anything, "d-1", "u-1").
			Return("url", nil)

		req := httptest.NewRequest("GET", "/draft/d-1/image", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	}
}
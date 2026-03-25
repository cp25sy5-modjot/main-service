package txihandler_test

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
	txhandler "github.com/cp25sy5-modjot/main-service/internal/transaction_item/handler"
	"github.com/cp25sy5-modjot/main-service/internal/transaction_item/mocks"
)

func TestTransactionItemHandler(t *testing.T) {

	jwt.GetUserIDFromClaims = func(c *fiber.Ctx) (string, error) {
		return "u-1", nil
	}

	// ===================== GET BY ID =====================
	t.Run("GetByID Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Get("/transactions/:transaction_id/item/:item_id", h.GetByID)

		mockSvc.EXPECT().
			GetByID(mock.Anything).
			Return(&e.TransactionItem{
				TransactionID: "tx-1",
				ItemID:        "item-1",
				Title:         "Food",
				Price:         100,
				CategoryID:    "c-1",
				Category: e.Category{
					CategoryName: "Food",
					ColorCode:    "#fff",
				},
			}, nil)

		req := httptest.NewRequest("GET", "/transactions/tx-1/item/item-1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("GetByID NotFound", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Get("/transactions/:transaction_id/item/:item_id", h.GetByID)

		mockSvc.EXPECT().
			GetByID(mock.Anything).
			Return(nil, errors.New("not found"))

		req := httptest.NewRequest("GET", "/transactions/tx-1/item/item-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 404, resp.StatusCode)
	})

	// ===================== UPDATE =====================
	t.Run("Update Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Put("/transactions/:transaction_id/item/:item_id", h.Update)

		body := `{
			"title": "New Title",
			"price": 200,
			"category_id": "c-2"
		}`

		mockSvc.EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(&e.TransactionItem{
				TransactionID: "tx-1",
				ItemID:        "item-1",
				Title:         "New Title",
				Price:         200,
				CategoryID:    "c-2",
				Category: e.Category{
					CategoryName: "New",
					ColorCode:    "#000",
				},
			}, nil)

		req := httptest.NewRequest(
			"PUT",
			"/transactions/tx-1/item/item-1",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Update Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Put("/transactions/:transaction_id/item/:item_id", h.Update)

		body := `{
			"title": "New Title",
			"price": 200,
			"category_id": "c-2"
		}`

		mockSvc.EXPECT().
			Update(mock.Anything, mock.Anything).
			Return(nil, errors.New("update failed"))

		req := httptest.NewRequest(
			"PUT",
			"/transactions/tx-1/item/item-1",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})

	// ===================== DELETE =====================
	t.Run("Delete Success", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Delete("/transactions/:transaction_id/item/:item_id", h.Delete)

		mockSvc.EXPECT().
			Delete(mock.Anything).
			Return(nil)

		req := httptest.NewRequest("DELETE", "/transactions/tx-1/item/item-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("Delete Error", func(t *testing.T) {
		app := fiber.New()
		mockSvc := mocks.NewMockService(t)

		h := txhandler.NewHandler(mockSvc)
		app.Delete("/transactions/:transaction_id/item/:item_id", h.Delete)

		mockSvc.EXPECT().
			Delete(mock.Anything).
			Return(errors.New("delete failed"))

		req := httptest.NewRequest("DELETE", "/transactions/tx-1/item/item-1", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, 500, resp.StatusCode)
	})
}
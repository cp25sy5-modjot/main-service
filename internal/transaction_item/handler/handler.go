package transactionhandler

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	txisvc "github.com/cp25sy5-modjot/main-service/internal/transaction_item/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service txisvc.Service
}

func NewHandler(svc txisvc.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

// GET /transactions/:transaction_id/item/:item_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	TransactionItemSearchParams, err := createTransactionItemSearchParams(c)
	if err != nil {
		return err
	}

	resp, err := h.service.GetByID(TransactionItemSearchParams)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction item not found")
	}

	return sresp.OK(c, buildTransactionItemResponse(resp), "Transaction item retrieved successfully")
}

// PUT /transactions/:transaction_id/item/:item_id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req m.TransactionItemUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	TransactionItemSearchParams, err := createTransactionItemSearchParams(c)
	if err != nil {
		return err
	}

	input := parseTransactionUpdateReqToServiceInput(&req)

	resp, err := h.service.Update(TransactionItemSearchParams, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update transaction")
	}
	return sresp.OK(c, buildTransactionItemResponse(resp), "Transaction item updated successfully")
}

// DELETE /transactions/:transaction_id/item/:item_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	TransactionItemSearchParams, err := createTransactionItemSearchParams(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(TransactionItemSearchParams); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete transaction item")
	}

	return sresp.OK(c, nil, "Transaction item deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	item_id := c.Params("item_id")
	if tx_id == "" || item_id == "" {
		return "", "", fiber.NewError(fiber.StatusBadRequest, "transaction_id and item_id parameters are required")
	}
	return tx_id, item_id, nil
}

func createTransactionItemSearchParams(c *fiber.Ctx) (*m.TransactionItemSearchParams, error) {
	tx_id, item_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &m.TransactionItemSearchParams{
		TransactionID: tx_id,
		ItemID:        item_id,
		UserID:        userID,
	}, nil
}

func buildTransactionItemResponse(item *e.TransactionItem) *m.TransactionItemRes {
	return &m.TransactionItemRes{
		TransactionID:     item.TransactionID,
		ItemID:            item.ItemID,
		Title:             item.Title,
		Price:             item.Price,
		CategoryID:        item.CategoryID,
		CategoryName:      item.Category.CategoryName,
		CategoryColor: item.Category.ColorCode,
	}
}

func parseTransactionUpdateReqToServiceInput(
	req *m.TransactionItemUpdateReq,
) *txisvc.TransactionItemUpdateInput {
	return &txisvc.TransactionItemUpdateInput{
		Title:      req.Title,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}
}

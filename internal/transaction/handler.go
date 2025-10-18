package transaction

import (
	"github.com/cp25sy5-modjot/main-service/internal/auth"
	successResp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// POST /transactions
func (h *Handler) Create(c *fiber.Ctx) error {
	var req TransactionInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	var tx Transaction
	_ = utils.MapStructs(&req, &tx)
	userID, err := auth.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	tx.UserID = userID
	tx.Type = "manual"
	if err := h.service.Create(&tx); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return successResp.Created(c, nil, "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := auth.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	transactions, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve transactions")
	}
	var resp []TransactionRes
	_ = utils.MapStructs(&transactions, &resp)
	return successResp.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id/product/:product_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	searchParams, err := createSearchParams(c)
	if err != nil {
		return err
	}
	transaction, err := h.service.GetByID(searchParams)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Transaction not found")
	}
	var resp TransactionRes
	_ = utils.MapStructs(&transaction, &resp)
	return successResp.OK(c, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:transaction_id/product/:product_id
func (h *Handler) Update(c *fiber.Ctx) error {
	var req TransactionUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	searchParams, err := createSearchParams(c)
	if err != nil {
		return err
	}
	if err := h.service.Update(searchParams, &req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update transaction")
	}
	var resp TransactionRes
	_ = utils.MapStructs(&req, &resp)
	return successResp.OK(c, resp, "Transaction updated successfully")
}

// DELETE /transactions/:transaction_id/product/:product_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	searchParams, err := createSearchParams(c)
	if err != nil {
		return err
	}
	if err := h.service.Delete(searchParams); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete transaction")
	}
	return successResp.OK(c, nil, "Transaction deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	prod_id := c.Params("product_id")
	if tx_id == "" || prod_id == "" {
		return "", "", fiber.NewError(fiber.StatusBadRequest, "transaction_id and product_id parameters are required")
	}
	return tx_id, prod_id, nil
}

func createSearchParams(c *fiber.Ctx) (*SearchParams, error) {
	tx_id, prod_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := auth.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &SearchParams{
		TransactionID: tx_id,
		ProductID:     prod_id,
		UserID:        userID,
	}, nil
}

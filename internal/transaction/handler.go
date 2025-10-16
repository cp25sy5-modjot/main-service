package transaction

import (
	r "modjot/internal/response"
	"modjot/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
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
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body", err)
	}
	// validate struct
	if err := r.Validator().Struct(req); err != nil {
		return r.UnprocessableEntity(c, "Validation Failed", err)

	}
	var tx Transaction
	_ = copier.Copy(&tx, &req)
	userID, err := utils.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	tx.UserID = userID
	tx.Type = "manual"
	if err := h.service.Create(&tx); err != nil {
		return r.InternalServerError(c, "Failed to create transaction")
	}

	return r.Created(c, nil, "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	userID, err := utils.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	transactions, err := h.service.GetAllByUserID(userID)
	if err != nil {
		return r.InternalServerError(c, "Failed to retrieve transactions")
	}
	var resp []TransactionRes
	_ = copier.Copy(&resp, &transactions)
	return r.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id/product/:product_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	searchParams, err := createSearchParams(c)
	if err != nil {
		return err
	}
	transaction, err := h.service.GetByID(searchParams)
	if err != nil {
		return r.NotFound(c, "Transaction not found")
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &transaction)
	return r.OK(c, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:transaction_id/product/:product_id
func (h *Handler) Update(c *fiber.Ctx) error {
	tx_id, prod_id, err := getTxIDAndProdID(c)
	if err != nil {
		return err
	}
	userID, err := utils.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req Transaction
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body", err)
	}
	req.TransactionID = tx_id
	req.ProductID = prod_id
	req.UserID = userID

	if err := h.service.Update(&req); err != nil {
		return r.InternalServerError(c, "Failed to update transaction")
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &req)
	return r.OK(c, resp, "Transaction updated successfully")
}

// DELETE /transactions/:transaction_id/product/:product_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	searchParams, err := createSearchParams(c)
	if err != nil {
		return err
	}
	if err := h.service.Delete(searchParams); err != nil {
		return r.InternalServerError(c, "Failed to delete transaction")
	}
	return r.OK(c, nil, "Transaction deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	prod_id := c.Params("product_id")
	if tx_id == "" || prod_id == "" {
		return "", "", r.BadRequest(c, "transaction_id and product_id parameters are required", nil)
	}
	return tx_id, prod_id, nil
}

func createSearchParams(c *fiber.Ctx) (*SearchParams, error) {
	tx_id, prod_id, err := getTxIDAndProdID(c)
	if err != nil {
		return nil, err
	}
	userID, err := utils.GetUserIDFromClaims(c)
	if err != nil {
		return nil, err
	}
	return &SearchParams{
		TransactionID: tx_id,
		ProductID:     prod_id,
		UserID:        userID,
	}, nil
}

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
		return r.BadRequest(c, "Invalid JSON body")
	}
	// validate struct
	if err := r.Validator().Struct(req); err != nil {
		return r.UnprocessableEntity(c, "Validation Failed", r.MapValidationErrors(err)...)

	}
	var tx Transaction
	_ = copier.Copy(&tx, &req)
	userID, err := getUserID(c)
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
	userID, err := getUserID(c)
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
	transaction, err := h.service.GetByID(createSearchParams(c))
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
	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	var req Transaction
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
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
	if err := h.service.Delete(createSearchParams(c)); err != nil {
		return r.InternalServerError(c, "Failed to delete transaction")
	}
	return r.OK(c, nil, "Transaction deleted successfully")
}

// utils
func getTxIDAndProdID(c *fiber.Ctx) (string, string, error) {
	tx_id := c.Params("transaction_id")
	prod_id := c.Params("product_id")
	if tx_id == "" || prod_id == "" {
		return "", "", fiber.NewError(fiber.StatusBadRequest, "ID parameter is required")
	}
	return tx_id, prod_id, nil
}

func getUserID(c *fiber.Ctx) (string, error) {
	userID := utils.GetUserIDFromClaims(c)
	if userID == "" {
		return "", r.InternalServerError(c, "Failed to get user ID from claims")
	}
	return userID, nil
}

func createSearchParams(c *fiber.Ctx) *SearchParams {
	tx_id, prod_id, _ := getTxIDAndProdID(c)
	userID, _ := getUserID(c)
	return &SearchParams{
		TransactionID: tx_id,
		ProductID:     prod_id,
		UserID:        userID,
	}
}

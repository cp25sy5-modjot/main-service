package transaction

import (
	"strconv"

	r "modjot/internal/response"

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
	if err := h.service.Create(&req); err != nil {
		return r.InternalServerError(c, "Failed to create transaction")
	}

	return r.Created(c, nil, "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	transactions, err := h.service.GetAll()
	if err != nil {
		return r.InternalServerError(c, "Failed to retrieve transactions")
	}
	var resp []TransactionRes
	_ = copier.Copy(&resp, &transactions)
	return r.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:transaction_id/product/:product_id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	tx_id := c.Params("transaction_id")
	prod_id := c.Params("")
	if tx_id == "" || prod_id == "" {
		return r.BadRequest(c, "ID parameter is required")
	}
	transaction, err := h.service.GetByID(tx_id, prod_id)
	if err != nil {
		return r.NotFound(c, "Transaction not found")
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &transaction)
	return r.OK(c, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:transaction_id/product/:product_id
func (h *Handler) Update(c *fiber.Ctx) error {

	tx_id := c.Params("transaction_id")
	prod_id := c.Params("product_id")

	if tx_id == "" || prod_id == "" {
		return r.BadRequest(c, "ID parameter is required")
	}

	var req Transaction
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
	}
	req.TransactionID = tx_id
	req.ProductID = prod_id

	if err := h.service.Update(&req); err != nil {
		return r.InternalServerError(c, "Failed to update transaction")
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &req)
	return r.OK(c, resp, "Transaction updated successfully")
}

// DELETE /transactions/:transaction_id/product/:product_id
func (h *Handler) Delete(c *fiber.Ctx) error {
	tx_id := c.Params("transaction_id")
	prod_id := c.Params("product_id")
	if tx_id == "" || prod_id == "" {
		return r.BadRequest(c, "ID parameter is required")
	}
	if err := h.service.Delete(tx_id, prod_id); err != nil {
		return r.InternalServerError(c, "Failed to delete transaction: "+strconv.Itoa(err.(*fiber.Error).Code))
	}
	return r.OK(c, nil, "Transaction deleted successfully")
}

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
	var req TransactionReq
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
	}

	// validate struct
	if err := r.Validator().Struct(req); err != nil {
		return r.UnprocessableEntity(c, "Validation Failed", r.MapValidationErrors(err)...)

	}
	if err := h.service.Create(&req); err != nil {
		return r.InternalError(c, "Failed to create transaction")
	}

	return r.Created(c, nil, "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	transactions, err := h.service.GetAll()
	if err != nil {
		return r.InternalError(c, err.Error())
	}
	var resp []TransactionRes
	_ = copier.Copy(&resp, &transactions)
	return r.OK(c, resp, "Transactions retrieved successfully")
}

// GET /transactions/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return r.BadRequest(c, err.Error())
	}
	transaction, err := h.service.GetByID(uint(id))
	if err != nil {
		return r.NotFound(c, err.Error())
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &transaction)
	return r.OK(c, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return r.BadRequest(c, err.Error())
	}

	var req Transaction
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, err.Error())
	}
	req.ID = uint(id)

	if err := h.service.Update(&req); err != nil {
		return r.InternalError(c, err.Error())
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &req)
	return r.OK(c, resp, "Transaction updated successfully")
}

// DELETE /transactions/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return r.BadRequest(c, err.Error())
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return r.InternalError(c, err.Error())
	}
	return r.OK(c, nil, "Transaction deleted successfully")
}

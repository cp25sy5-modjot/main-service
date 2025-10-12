package transaction

import (
	"strconv"

	mw "modjot/internal/middleware"

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
		return mw.WriteError(c, fiber.StatusBadRequest, "Invalid request", "bad_request", err.Error(), nil)
	}
	if err := h.service.Create(&req); err != nil {
		return mw.WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", err.Error(), nil)
	}

	return mw.WriteSuccess(c, fiber.StatusCreated, nil, "Transaction created successfully")
}

// GET /transactions
func (h *Handler) GetAll(c *fiber.Ctx) error {
	transactions, err := h.service.GetAll()
	if err != nil {
		return mw.WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", err.Error(), nil)
	}
	var resp []TransactionRes
	_ = copier.Copy(&resp, &transactions)
	return mw.WriteSuccess(c, fiber.StatusOK, resp, "Transactions retrieved successfully")
}

// GET /transactions/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return mw.WriteError(c, fiber.StatusBadRequest, "Invalid ID", "bad_request", err.Error(), nil)
	}
	transaction, err := h.service.GetByID(uint(id))
	if err != nil {
		return mw.WriteError(c, fiber.StatusNotFound, "Transaction not found", "not_found", err.Error(), nil)
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &transaction)
	return mw.WriteSuccess(c, fiber.StatusOK, resp, "Transaction retrieved successfully")
}

// PUT /transactions/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return mw.WriteError(c, fiber.StatusBadRequest, "Invalid ID", "bad_request", err.Error(), nil)
	}

	var req Transaction
	if err := c.BodyParser(&req); err != nil {
		return mw.WriteError(c, fiber.StatusBadRequest, "Invalid request", "bad_request", err.Error(), nil)
	}
	req.ID = uint(id)

	if err := h.service.Update(&req); err != nil {
		return mw.WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", err.Error(), nil)
	}
	var resp TransactionRes
	_ = copier.Copy(&resp, &req)
	return mw.WriteSuccess(c, fiber.StatusOK, resp, "Transaction updated successfully")
}

// DELETE /transactions/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return mw.WriteError(c, fiber.StatusBadRequest, "Invalid ID", "bad_request", err.Error(), nil)
	}
	if err := h.service.Delete(uint(id)); err != nil {
		return mw.WriteError(c, fiber.StatusInternalServerError, "Internal Server Error", "internal_error", err.Error(), nil)
	}
	return mw.WriteSuccess(c, fiber.StatusNoContent, nil, "Transaction deleted successfully")
}

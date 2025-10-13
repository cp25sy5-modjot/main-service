package user

import (
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

// POST /users
func (h *Handler) Create(c *fiber.Ctx) error {
	var req UserInsertReq
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
	}

	// validate struct
	if err := r.Validator().Struct(req); err != nil {
		return r.UnprocessableEntity(c, "Validation Failed", r.MapValidationErrors(err)...)
	}

	if err := h.service.Create(&req); err != nil {
		return r.InternalServerError(c, err.Error())
	}

	return r.Created(c, nil, "User created successfully")
}

// PUT /users/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return r.BadRequest(c, "ID parameter is required")
	}

	var req UserUpdateReq
	if err := c.BodyParser(&req); err != nil {
		return r.BadRequest(c, "Invalid JSON body")
	}

	// validate struct
	if err := r.Validator().Struct(req); err != nil {
		return r.UnprocessableEntity(c, "Validation Failed", r.MapValidationErrors(err)...)
	}
	var entity User
	_ = copier.Copy(&entity, &req)
	entity.UserID = id

	if err := h.service.Update(&entity); err != nil {
		return r.InternalServerError(c, err.Error())
	}

	return r.OK(c, nil, "User updated successfully")
}

// DELETE /users/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return r.BadRequest(c, "ID parameter is required")
	}
	if err := h.service.Delete(id); err != nil {
		return r.InternalServerError(c, err.Error())
	}
	return r.OK(c, nil, "User deleted successfully")
}

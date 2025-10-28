package user

import (
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
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

// POST /users
func (h *Handler) Create(c *fiber.Ctx) error {
	var req UserInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.service.Create(&req); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return successResp.Created(c, nil, "User created successfully")
}

// PUT /users/:id
func (h *Handler) Update(c *fiber.Ctx) error {

	var req UserUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var entity User
	utils.MapNonNilStructs(&req, &entity)
	entity.UserID = userID

	if err := h.service.Update(&entity); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return successResp.OK(c, nil, "User updated successfully")
}

// DELETE /users/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	if err := h.service.Delete(userID); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return successResp.OK(c, nil, "User deleted successfully")
}

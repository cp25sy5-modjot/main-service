package user

import (
	"log"

	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	successResp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	model "github.com/cp25sy5-modjot/main-service/internal/user/model"
	svc "github.com/cp25sy5-modjot/main-service/internal/user/service"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *svc.Service
}

func NewHandler(service *svc.Service) *Handler {
	return &Handler{service}
}

// GET /user
func (h *Handler) GetSelf(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	user, err := h.service.GetByID(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	//parse to response model
	var userRes model.UserRes
	utils.MapStructs(user, &userRes)
	return successResp.OK(c, userRes, "User retrieved successfully")
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req model.UserInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.service.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var userRes model.UserRes
	utils.MapStructs(user, &userRes)
	return successResp.Created(c, userRes, "User created successfully")
}

// PUT /user
func (h *Handler) Update(c *fiber.Ctx) error {

	var req model.UserUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	updated, err := h.service.Update(userID, &req);
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var resp model.UserRes
	utils.MapStructs(updated, &resp)

	return successResp.OK(c, resp, "User updated successfully")
}

// DELETE /user
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

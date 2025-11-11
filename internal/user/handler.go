package user

import (
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	successResp "github.com/cp25sy5-modjot/main-service/internal/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
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
	var userRes UserRes
	utils.MapStructs(&user, &userRes)
	log.Printf("User retrieved: %+v", userRes)
	return successResp.OK(c, userRes, "User retrieved successfully")
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req UserInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.service.Create(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	var userRes UserRes
	utils.MapNonNilStructs(user, &userRes)
	return successResp.Created(c, userRes, "User created successfully")
}

// PUT /user
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

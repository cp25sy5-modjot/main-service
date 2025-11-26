package userhandler

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	sresp "github.com/cp25sy5-modjot/main-service/internal/response/success"
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

	return sresp.OK(c, buildUserResponse(user), "User retrieved successfully")
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var req m.UserInsertReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	input := &svc.UserCreateInput{
		UserBinding: e.UserBinding{
			GoogleID:   req.UserBinding.GoogleID,
			FacebookID: req.UserBinding.FacebookID,
			AppleID:    req.UserBinding.AppleID,
		},
		Name: req.Name,
		DOB:  utils.NormalizeToUTC(req.DOB, ""),
	}

	user, err := h.service.Create(input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sresp.Created(c, buildUserResponse(user), "User created successfully")
}

// PUT /user
func (h *Handler) Update(c *fiber.Ctx) error {

	var req m.UserUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}
	input := &svc.UserUpdateInput{
		Name: req.Name,
		DOB:  utils.NormalizeToUTC(req.DOB, ""),
	}

	updated, err := h.service.Update(userID, input)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return sresp.OK(c, buildUserResponse(updated), "User updated successfully")
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
	return sresp.OK(c, nil, "User deleted successfully")
}

func buildUserResponse(user *e.User) *m.UserRes {
	return &m.UserRes{
		Name:      user.Name,
		DOB:       utils.ToUserLocal(user.DOB, ""),
		Status:    string(user.Status),
		CreatedAt: utils.ToUserLocal(user.CreatedAt, ""),
		UserBinding: m.UserBinding{
			GoogleID:   user.UserBinding.GoogleID,
			FacebookID: user.UserBinding.FacebookID,
			AppleID:    user.UserBinding.AppleID,
		},
	}
}

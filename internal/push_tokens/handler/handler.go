package pushhandler

import (
	"github.com/gofiber/fiber/v2"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	pushsvc "github.com/cp25sy5-modjot/main-service/internal/push_tokens/service"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
)

type Handler struct {
	service pushsvc.Service
}

func NewHandler(svc pushsvc.Service) *Handler {
	return &Handler{service: svc}
}

// POST /push/register
func (h *Handler) Register(c *fiber.Ctx) error {
	var req m.PushRegisterReq

	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}

	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	err = h.service.Register(c.Context(), userID, req.Token, req.Platform)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to register push token")
	}

	return sresp.OK(c, nil, "Push token registered successfully")
}

// POST /push/test
func (h *Handler) TestSend(ctx *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(ctx)
	if err != nil {
		return err
	}

	err = h.service.Send(ctx.Context(), userID, "Test Notification", "This is a test push notification")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to send test notification")
	}

	return sresp.OK(ctx, nil, "Test notification sent successfully")
}

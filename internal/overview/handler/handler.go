package overviewhandler

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	overviewsvc "github.com/cp25sy5-modjot/main-service/internal/overview/service"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
)

type Handler struct {
	service overviewsvc.Service
}

func NewHandler(s overviewsvc.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetOverview(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	dateStr := c.Query("date")
	var base time.Time

	if dateStr == "" {
		base = time.Now().UTC()
	} else {
		base, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		}
	}

	res, err := h.service.GetOverview(userID, base)
	if err != nil {
		return err
	}

	return sresp.OK(c, res, "Overview retrieved successfully")
}

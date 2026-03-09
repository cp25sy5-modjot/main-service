package fixcosthandler

import (
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	fcsvc "github.com/cp25sy5-modjot/main-service/internal/fix_cost/service"
	"github.com/cp25sy5-modjot/main-service/internal/jwt"
	"github.com/cp25sy5-modjot/main-service/internal/mapper"
	sresp "github.com/cp25sy5-modjot/main-service/internal/shared/response/success"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service fcsvc.Service
}

func NewHandler(service fcsvc.Service) *Handler {
	return &Handler{service}
}

// POST /fix_cost
func (h *Handler) Create(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	var req m.FixCostCreateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	input := mapper.ParseFixCostCreateReqToServiceInput(userID, &req)

	fixCost, err := h.service.Create(c.Context(), input)
	if err != nil {
		return err
	}

	return sresp.Created(c, mapper.BuildFixCostResponse(fixCost), "Fix cost created successfully")
}

// PUT /fix_cost/:id
func (h *Handler) Update(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	fixCostID := c.Params("id")
	if fixCostID == "" {
		return fiber.ErrBadRequest
	}

	var req m.FixCostUpdateReq
	if err := utils.ParseBodyAndValidate(c, &req); err != nil {
		return err
	}
	input := mapper.ParseFixCostUpdateReqToServiceInput(userID, fixCostID, &req)

	if err := h.service.Update(c.Context(), input); err != nil {
		return err
	}

	return sresp.NoContent(c)
}

// DELETE /fix_cost/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	fixCostID := c.Params("id")
	if fixCostID == "" {
		return fiber.ErrBadRequest
	}

	if err := h.service.Delete(c.Context(), fixCostID, userID); err != nil {
		return err
	}

	return sresp.NoContent(c)
}

// GET /fix_cost/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	fixCostID := c.Params("id")
	if fixCostID == "" {
		return fiber.ErrBadRequest
	}

	fixCost, err := h.service.GetByID(c.Context(), fixCostID, userID)
	if err != nil {
		return err
	}

	return sresp.OK(c, mapper.BuildFixCostResponse(fixCost), "Fix cost retrieved successfully")
}

// GET /fix_cost
func (h *Handler) GetAllByUserID(c *fiber.Ctx) error {
	userID, err := jwt.GetUserIDFromClaims(c)
	if err != nil {
		return err
	}

	fixCosts, err := h.service.GetAllByUserID(c.Context(), userID)
	if err != nil {
		return err
	}

	var res []*m.FixCostRes
	for _, fc := range fixCosts {
		res = append(res, mapper.BuildFixCostResponse(fc))
	}
	return sresp.OK(c, res, "Fix costs retrieved successfully")
}

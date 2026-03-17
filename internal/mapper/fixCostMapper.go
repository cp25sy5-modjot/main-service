package mapper

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

func ParseFixCostCreateReqToServiceInput(
	userID string,
	req *m.FixCostCreateReq,
) *m.FixCostCreateInput {
	return &m.FixCostCreateInput{
		UserID:        userID,
		Title:         req.Title,
		Price:         req.Price,
		CategoryID:    req.CategoryID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		RemainingRuns: req.RemainingRuns,
		IntervalType:  req.IntervalType,
		IntervalValue: req.IntervalValue,
	}
}

func ParseFixCostUpdateReqToServiceInput(
	userID string,
	fixCostID string,
	req *m.FixCostUpdateReq,
) *m.FixCostUpdateInput {
	return &m.FixCostUpdateInput{
		UserID:        userID,
		FixCostID:     fixCostID,
		Title:         req.Title,
		Price:         req.Price,
		CategoryID:    req.CategoryID,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		RemainingRuns: req.RemainingRuns,
		Status:        req.Status,
		IntervalType:  req.IntervalType,
		IntervalValue: req.IntervalValue,
	}
}

func BuildFixCostResponse(fc *e.FixCost) *m.FixCostRes {
	return &m.FixCostRes{
		FixCostID:     fc.FixCostID,
		Title:         fc.Title,
		Price:         fc.Price,
		CategoryID:    fc.CategoryID,
		StartDate:     fc.StartDate,
		EndDate:       fc.EndDate,
		RemainingRuns: fc.RemainingRuns,
		IntervalType:  string(fc.IntervalType),
		IntervalValue: fc.IntervalValue,
		Status:        string(fc.Status),
		NextRunDate:   fc.NextRunDate,

		CategoryIcon:  fc.Category.Icon,
		CategoryColor: fc.Category.ColorCode,
		CategoryName:  fc.Category.CategoryName,

		CreatedAt: fc.CreatedAt,
		UpdatedAt: fc.UpdatedAt,
	}
}

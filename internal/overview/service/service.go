package overviewsvc

import (
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	overviewrepo "github.com/cp25sy5-modjot/main-service/internal/overview/repository"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
)

// Service defines behavior for overview use case
type Service interface {
	// GetOverview returns:
	// - last N transactions (limit internal)
	// - top categories by spending (limit internal, current month from baseDate)
	GetOverview(userID string, baseDate time.Time) (*m.OverviewResponse, error)
}

type service struct {
	repo *overviewrepo.Repository
}

func NewService(repo *overviewrepo.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetOverview(userID string, t time.Time) (*m.OverviewResponse, error) {
	// Normalize baseDate just to be safe; use its location for month range
	start, end := utils.GetStartAndEndOfMonth(t)
	// 1) last 3 transactions (global, not month-limited)
	lastTx, err := s.repo.GetLastTransactions(userID, start, end, 3)
	if err != nil {
		return nil, err
	}

	// 2) top 3 categories by spending in that month
	topCats, err := s.repo.GetTopCategoriesBySpending(userID, start, end, 3)
	if err != nil {
		return nil, err
	}

	// 3) current month total
	currentMonthTotal, err := s.repo.GetMonthTotal(userID, start, end)
	if err != nil {
		return nil, err
	}

	// 4) previous month total
	previousMonthTotal, err := s.repo.GetMonthTotal(userID, start.AddDate(0, -1, 0), end.AddDate(0, -1, 0))
	if err != nil {
		return nil, err
	}

	// 5) build response
	return &m.OverviewResponse{
		LastTransactions:   lastTx,
		TopCategories:      topCats,
		CurrentMonthTotal:  currentMonthTotal,
		PreviousMonthTotal: previousMonthTotal,
	}, nil
}

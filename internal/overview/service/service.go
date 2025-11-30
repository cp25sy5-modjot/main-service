package overviewsvc

import (
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	overviewrepo "github.com/cp25sy5-modjot/main-service/internal/overview/repository"
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

func (s *service) GetOverview(userID string, baseDate time.Time) (*m.OverviewResponse, error) {
	// Normalize baseDate just to be safe; use its location for month range
	startOfMonth := time.Date(baseDate.Year(), baseDate.Month(), 1, 7, 0, 0, 0, baseDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0) // first day of next month (exclusive)

	// 1) last 3 transactions (global, not month-limited)
	lastTx, err := s.repo.GetLastTransactions(userID, 3)
	if err != nil {
		return nil, err
	}

	// 2) top 3 categories by spending in that month
	topCats, err := s.repo.GetTopCategoriesBySpending(userID, startOfMonth, endOfMonth, 3)
	if err != nil {
		return nil, err
	}

	// 3) build response
	return &m.OverviewResponse{
		LastTransactions: lastTx,
		TopCategories:    topCats,
	}, nil
}

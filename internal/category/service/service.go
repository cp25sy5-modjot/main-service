package categorysvc

import (
	"time"

	categoryrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	utils "github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	"github.com/google/uuid"
)

type Service interface {
	Create(userId string, input *CategoryCreateInput) (*e.Category, error)
	GetAllByUserID(userID string) ([]e.Category, error)
	GetAllByUserIDWithTransactions(userID string, filter *m.TransactionFilter) ([]e.Category, error)
	GetByID(params *m.CategorySearchParams) (*e.Category, error)
	GetByIDWithTransactions(params *m.CategorySearchParams, filter *m.TransactionFilter) (*e.Category, error)
	Update(params *m.CategorySearchParams, input *CategoryUpdateInput) (*e.Category, error)
	Delete(params *m.CategorySearchParams) error
}

type service struct {
	categoryrepo *categoryrepo.Repository
	txrepo       *txrepo.Repository
}

func NewService(categoryrepo *categoryrepo.Repository, txrepo *txrepo.Repository) *service {
	return &service{categoryrepo, txrepo}
}

func (s *service) Create(userId string, input *CategoryCreateInput) (*e.Category, error) {
	cate := &e.Category{
		CategoryID:   uuid.New().String(),
		CategoryName: input.CategoryName,
		UserID:       userId,
		Budget:       input.Budget,
		ColorCode:    input.ColorCode,
		CreatedAt:    time.Now(),
	}
	return saveNewCategory(s, cate)
}

func (s *service) GetAllByUserID(userID string) ([]e.Category, error) {
	categories, err := s.categoryrepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetAllByUserIDWithTransactions(userID string, filter *m.TransactionFilter) ([]e.Category, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	start, end := getMonthRange(filter)

	categories, err := s.categoryrepo.FindAllByUserIDWithTransactionsFiltered(userID, start, end)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *service) GetByID(params *m.CategorySearchParams) (*e.Category, error) {
	category, err := s.categoryrepo.FindByID(params)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) GetByIDWithTransactions(params *m.CategorySearchParams, filter *m.TransactionFilter) (*e.Category, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	start, end := getMonthRange(filter)

	category, err := s.categoryrepo.FindByIDWithTransactionsFiltered(params, start, end)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *service) Update(params *m.CategorySearchParams, input *CategoryUpdateInput) (*e.Category, error) {
	exists, err := s.categoryrepo.FindByID(params)
	if err != nil {
		return nil, err
	}

	if err := utils.MapStructs(input, exists); err != nil {
		return nil, err
	}

	updatedCat, err := s.categoryrepo.Update(exists)
	if err != nil {
		return nil, err
	}

	return updatedCat, nil
}

func (s *service) Delete(params *m.CategorySearchParams) error {
	_, err := s.categoryrepo.FindByID(params)
	if err != nil {
		return err
	}
	return s.categoryrepo.Delete(params)
}

// utils functions for service

func applyUpdates(cat *e.Category, in *CategoryUpdateInput) {
	cat.CategoryName = in.CategoryName
	cat.Budget = in.Budget
	cat.ColorCode = in.ColorCode
}

func saveNewCategory(s *service, cat *e.Category) (*e.Category, error) {
	newCat, err := s.categoryrepo.Create(cat)
	if err != nil {
		return nil, err
	}
	// Reload with preload
	catWithDetails, err := s.categoryrepo.FindByID(&m.CategorySearchParams{
		CategoryID: &newCat.CategoryID,
		UserID:     newCat.UserID,
	})
	if err != nil {
		return nil, err
	}
	return catWithDetails, nil
}

func getMonthRange(filter *m.TransactionFilter) (time.Time, time.Time) {
	t := *filter.Date
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 7, 0, 0, 0, t.Location())

	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return startOfMonth, endOfMonth
}

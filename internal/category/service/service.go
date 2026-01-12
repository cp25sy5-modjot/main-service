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
	CreateDefaultCategories(userID string) error
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
		CreatedAt:    time.Now().UTC(),
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
		now := time.Now().UTC()
		filter.Date = &now
	}
	start, end := utils.GetStartAndEndOfMonth(*filter.Date)

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
		now := time.Now().UTC()
		filter.Date = &now
	}
	start, end := utils.GetStartAndEndOfMonth(*filter.Date)

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

func (s *service) CreateDefaultCategories(userID string) error {
	defaultCategories := []string{"อาหาร", "การเดินทาง", "ความบันเทิง", "ชอปปิ้ง", "อื่นๆ"}
	for _, categoryName := range defaultCategories {
		_, err := s.categoryrepo.Create(&e.Category{
			CategoryID:   uuid.New().String(),
			CategoryName: categoryName,
			UserID:       userID,
			Budget:       1000.0,
			ColorCode:    utils.GenerateRandomColor(),
			CreatedAt:    time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// utils functions for service
func saveNewCategory(s *service, cat *e.Category) (*e.Category, error) {
	newCat, err := s.categoryrepo.Create(cat)
	if err != nil {
		return nil, err
	}
	// Reload with preload
	catWithDetails, err := s.categoryrepo.FindByID(&m.CategorySearchParams{
		CategoryID: newCat.CategoryID,
		UserID:     newCat.UserID,
	})
	if err != nil {
		return nil, err
	}
	return catWithDetails, nil
}

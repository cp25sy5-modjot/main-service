package categorysvc

import (
	categoryrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	utils "github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	categoryrepo *categoryrepo.Repository
	txrepo       *txrepo.Repository
}

func NewService(categoryrepo *categoryrepo.Repository, txrepo *txrepo.Repository) *Service {
	return &Service{categoryrepo, txrepo}
}

func (s *Service) Create(userId string, input *CategoryCreateInput) (*e.Category, error) {
	cate := &e.Category{
		CategoryID:   uuid.New().String(),
		CategoryName: input.CategoryName,
		UserID:       userId,
		Budget:       input.Budget,
		ColorCode:    input.ColorCode,
		CreatedAt:    utils.NowUTC(),
	}
	return saveNewCategory(s, cate)
}

func (s *Service) GetAllByUserID(userID string) ([]e.Category, error) {
	categories, err := s.categoryrepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetAllByUserIDWithTransactions(userID string) ([]e.Category, error) {
	categories, err := s.categoryrepo.FindAllByUserIDWithTransactions(userID)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetByID(params *m.CategorySearchParams) (*e.Category, error) {
	category, err := s.categoryrepo.FindByID(params)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) GetByIDWithTransactions(params *m.CategorySearchParams) (*e.Category, error) {
	category, err := s.categoryrepo.FindByIDWithTransactions(params)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) Update(params *m.CategorySearchParams, input *CategoryUpdateInput) (*e.Category, error) {
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

func (s *Service) Delete(params *m.CategorySearchParams) error {
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

func saveNewCategory(s *Service, cat *e.Category) (*e.Category, error) {
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

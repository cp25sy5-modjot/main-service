package category

import (
	"time"

	repo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(category *e.Category) (*m.CategoryRes, error) {
	cate := &e.Category{
		CategoryID:   uuid.New().String(),
		CategoryName: category.CategoryName,
		UserID:       category.UserID,
		Budget:       category.Budget,
		ColorCode:    category.ColorCode,
		CreatedAt:    time.Now(),
	}
	return saveNewCategory(s, cate)
}

func (s *Service) GetAllByUserID(userID string) ([]m.CategoryRes, error) {
	categories, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	return buildCategoryResponses(categories), nil
}

func (s *Service) GetAllByUserIDWithTransactions(userID string) ([]m.CategoryRes, error) {
	categories, err := s.repo.FindAllByUserIDWithTransactions(userID)
	if err != nil {
		return nil, err
	}
	return buildCategoryResponses(categories), nil
}

func (s *Service) GetByID(params *m.CategorySearchParams) (*m.CategoryRes, error) {
	category, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	return buildCategoryResponse(category), nil
}

func (s *Service) GetByIDWithTransactions(params *m.CategorySearchParams) (*m.CategoryRes, error) {
	category, err := s.repo.FindByIDWithTransactions(params)
	if err != nil {
		return nil, err
	}
	return buildCategoryResponse(category), nil
}

func (s *Service) Update(params *m.CategorySearchParams, category *m.CategoryUpdateReq) (*m.CategoryRes, error) {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	
	if err = utils.MapStructs(category, &exists); err != nil {
		return nil, err
	}

	updatedCat, err := s.repo.Update(exists)
	if err != nil {
		return nil, err
	}

	return buildCategoryResponse(updatedCat), nil
}

func (s *Service) Delete(params *m.CategorySearchParams) error {
	_, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	return s.repo.Delete(params)
}

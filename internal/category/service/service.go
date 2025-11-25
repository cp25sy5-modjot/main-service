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

func (s *Service) Create(category *e.Category) (*e.Category, error) {
	cate := &e.Category{
		CategoryID:   uuid.New().String(),
		CategoryName: category.CategoryName,
		UserID:       category.UserID,
		Budget:       category.Budget,
		ColorCode:    category.ColorCode,
		CreatedAt:    time.Now(),
	}
	return s.repo.Create(cate)
}

func (s *Service) GetAllByUserID(userID string) ([]e.Category, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *Service) GetByID(params *m.CategorySearchParams) (*e.Category, error) {
	return s.repo.FindByID(params)
}

func (s *Service) Update(params *m.CategorySearchParams, category *m.CategoryUpdateReq) (*e.Category, error) {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	if err := validateCategoryOwnership(exists, params.UserID); err != nil {
		return nil, err
	}
	if err = utils.MapStructs(category, &exists); err != nil {
		return nil, err
	}

	return exists, s.repo.Update(exists)
}

func (s *Service) Delete(params *m.CategorySearchParams) error {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	if err := validateCategoryOwnership(exists, params.UserID); err != nil {
		return err
	}

	return s.repo.Delete(params)
}

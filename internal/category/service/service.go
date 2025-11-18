package category

import (
	"time"

	model "github.com/cp25sy5-modjot/main-service/internal/category/model"
	repo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	r "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *repo.Repository
}

func NewService(repo *repo.Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(category *model.Category) (*model.Category, error) {
	cate := &model.Category{
		CategoryID:   uuid.New().String(),
		CategoryName: category.CategoryName,
		UserID:       category.UserID,
		Budget:       category.Budget,
		ColorCode:    category.ColorCode,
		CreatedAt:    time.Now(),
	}
	return s.repo.Create(cate)
}

func (s *Service) GetAllByUserID(userID string) ([]model.Category, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *Service) GetByID(params *model.CategorySearchParams) (*model.Category, error) {
	return s.repo.FindByID(params)
}

func (s *Service) Update(params *model.CategorySearchParams, category *model.CategoryUpdateReq) (*model.Category, error) {
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

func (s *Service) Delete(params *model.CategorySearchParams) error {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	if err := validateCategoryOwnership(exists, params.UserID); err != nil {
		return err
	}

	return s.repo.Delete(params)
}

func validateCategoryOwnership(cate *model.Category, userID string) error {
	if cate.UserID != userID {
		return r.Conflict(nil, "You are not authorized to access this category")
	}
	return nil
}

package user

import (
	"time"

	catModel "github.com/cp25sy5-modjot/main-service/internal/category/model"
	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	userModel "github.com/cp25sy5-modjot/main-service/internal/user/model"
	userRepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *userRepo.Repository
	cat  *catSvc.Service
}

func NewService(repo *userRepo.Repository, cat *catSvc.Service) *Service {
	return &Service{repo, cat}
}

func (s *Service) Create(user *userModel.UserInsertReq) (*userModel.User, error) {
	UserID := uuid.New().String()
	u := &userModel.User{
		UserID: UserID,
		UserBinding: userModel.UserBinding{
			GoogleID:   user.UserBinding.GoogleID,
			FacebookID: user.UserBinding.FacebookID,
			AppleID:    user.UserBinding.AppleID,
		},
		Email:     user.Email,
		Name:      user.Name,
		DOB:       user.DOB,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	defaultCategories := []string{"อาหาร", "การเดินทาง", "ความบันเทิง", "ชอปปิ้ง", "อื่นๆ"}
	for _, categoryName := range defaultCategories {
		_, err := s.cat.Create(&catModel.Category{
			CategoryName: categoryName,
			UserID:       userCreated.UserID,
			Budget:       1000.0,
		})
		if err != nil {
			return nil, err
		}
	}
	return userCreated, nil
}

func (s *Service) GetAll() ([]*userModel.User, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByEmail(email string) (*userModel.User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) GetByID(user_id string) (*userModel.User, error) {
	return s.repo.FindByID(user_id)
}

func (s *Service) GetByGoogleID(google_id string) (*userModel.User, error) {
	return s.repo.FindByGoogleID(google_id)
}

func (s *Service) GetByFacebookID(facebook_id string) (*userModel.User, error) {
	return s.repo.FindByFacebookID(facebook_id)
}

func (s *Service) GetByAppleID(apple_id string) (*userModel.User, error) {
	return s.repo.FindByAppleID(apple_id)
}

func (s *Service) Update(userID string, req *userModel.UserUpdateReq) (*userModel.User, error) {
	exists, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if err := utils.MapStructs(req, exists); err != nil {
		return nil, err
	}
	return exists, s.repo.Update(exists)
}

func (s *Service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

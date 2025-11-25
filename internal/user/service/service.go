package usersvc

import (
	catsvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	"github.com/google/uuid"
)

type Service struct {
	repo *userrepo.Repository
	cat  *catsvc.Service
}

func NewService(repo *userrepo.Repository, cat *catsvc.Service) *Service {
	return &Service{repo, cat}
}

func (s *Service) Create(input *UserCreateInput) (*e.User, error) {
	UserID := uuid.New().String()
	u := buildUserObjectToCreate(UserID, input)
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	if err := createDefaultCategories(s, UserID); err != nil {
		return nil, err
	}

	return userCreated, nil
}

func (s *Service) CreateMockUser(input *UserCreateInput, uid string) (*e.User, error) {
	u := buildUserObjectToCreate(uid, input)
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	if err := createDefaultCategories(s, uid); err != nil {
		return nil, err
	}

	return userCreated, nil
}

func (s *Service) GetAll() ([]*e.User, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(user_id string) (*e.User, error) {
	return s.repo.FindByID(user_id)
}

func (s *Service) GetByGoogleID(google_id string) (*e.User, error) {
	return s.repo.FindByGoogleID(google_id)
}

func (s *Service) GetByFacebookID(facebook_id string) (*e.User, error) {
	return s.repo.FindByFacebookID(facebook_id)
}

func (s *Service) GetByAppleID(apple_id string) (*e.User, error) {
	return s.repo.FindByAppleID(apple_id)
}

func (s *Service) Update(userID string, input *UserUpdateInput) (*e.User, error) {
	exists, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if err := utils.MapStructs(input, exists); err != nil {
		return nil, err
	}

	if exists.Status == e.StatusPreActive {
		exists.Status = e.StatusActive
	}

	updatedUser, err := s.repo.Update(exists)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *Service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

// utils functions for service
func buildUserObjectToCreate(uid string, input *UserCreateInput) *e.User {
	return &e.User{
		UserID: uid,
		UserBinding: e.UserBinding{
			GoogleID:   input.UserBinding.GoogleID,
			FacebookID: input.UserBinding.FacebookID,
			AppleID:    input.UserBinding.AppleID,
		},
		Name:      input.Name,
		DOB:       input.DOB,
		CreatedAt: utils.NowUTC(),
		UpdatedAt: utils.NowUTC(),
	}
}

func createDefaultCategories(s *Service, uid string) error {
	defaultCategories := []string{"อาหาร", "การเดินทาง", "ความบันเทิง", "ชอปปิ้ง", "อื่นๆ"}
	for _, categoryName := range defaultCategories {
		_, err := s.cat.Create(uid, &catsvc.CategoryCreateInput{
			CategoryName: categoryName,
			Budget:       1000.0,
			ColorCode:    utils.GenerateRandomColor(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

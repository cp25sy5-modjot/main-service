package user

import (
	"time"

	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
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

func (s *Service) Create(user *m.UserInsertReq) (*e.User, error) {
	UserID := uuid.New().String()
	u := buildUserObjectToCreate(UserID, user)
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}

	if err := createDefaultCategories(s, UserID); err != nil {
		return nil, err
	}

	return userCreated, nil
}

func (s *Service) CreateMockUser(user *m.UserInsertReq, uid string) (*e.User, error) {
	u := buildUserObjectToCreate(uid, user)
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

func (s *Service) GetByEmail(email string) (*e.User, error) {
	return s.repo.FindByEmail(email)
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

func (s *Service) Update(userID string, req *m.UserUpdateReq) (*e.User, error) {
	exists, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if err := utils.MapStructs(req, exists); err != nil {
		return nil, err
	}
	if exists.Status == e.StatusPreActive {
		exists.Status = e.StatusActive
	}
	return exists, s.repo.Update(exists)
}

func (s *Service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

func buildUserObjectToCreate(uid string, user *m.UserInsertReq) *e.User {
	return &e.User{
		UserID: uid,
		UserBinding: e.UserBinding{
			GoogleID:   user.UserBinding.GoogleID,
			FacebookID: user.UserBinding.FacebookID,
			AppleID:    user.UserBinding.AppleID,
		},
		Name:      user.Name,
		DOB:       user.DOB,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createDefaultCategories(s *Service, uid string) error {
	defaultCategories := []string{"อาหาร", "การเดินทาง", "ความบันเทิง", "ชอปปิ้ง", "อื่นๆ"}
	for _, categoryName := range defaultCategories {
		_, err := s.cat.Create(&e.Category{
			CategoryName: categoryName,
			UserID:       uid,
			Budget:       1000.0,
			ColorCode:    utils.GenerateRandomColor(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

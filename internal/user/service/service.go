package usersvc

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	"github.com/google/uuid"
)

type Service interface {
	Create(input *UserCreateInput) (*e.User, error)
	CreateMockUser(input *UserCreateInput, uid string) (*e.User, error)

	GetAll() ([]*e.User, error)
	GetByID(user_id string) (*e.User, error)
	GetByGoogleID(google_id string) (*e.User, error)
	// GetByFacebookID(facebook_id string) (*e.User, error)
	// GetByAppleID(apple_id string) (*e.User, error)

	Update(userID string, input *UserUpdateInput) (*e.User, error)
	Delete(user_id string) error
}

type service struct {
	repo *userrepo.Repository
}

func NewService(repo *userrepo.Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(input *UserCreateInput) (*e.User, error) {
	UserID := uuid.New().String()
	u := buildUserObjectToCreate(UserID, input)
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}
	return userCreated, nil
}

func (s *service) CreateMockUser(input *UserCreateInput, uid string) (*e.User, error) {
	u := buildUserObjectToCreate(uid, input)
	userCreated, err := s.repo.Create(u)
	if err != nil {
		return nil, err
	}
	return userCreated, nil
}

func (s *service) GetAll() ([]*e.User, error) {
	users, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *service) GetByID(user_id string) (*e.User, error) {
	user, err := s.repo.FindByID(user_id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetByGoogleID(google_id string) (*e.User, error) {
	user, err := s.repo.FindByGoogleID(google_id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// func (s *service) GetByFacebookID(facebook_id string) (*e.User, error) {
// 	user, err := s.repo.FindByFacebookID(facebook_id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func (s *service) GetByAppleID(apple_id string) (*e.User, error) {
// 	user, err := s.repo.FindByAppleID(apple_id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

func (s *service) Update(userID string, input *UserUpdateInput) (*e.User, error) {
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

func (s *service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

// utils functions for service
func buildUserObjectToCreate(uid string, input *UserCreateInput) *e.User {
	return &e.User{
		UserID: uid,
		UserBinding: e.UserBinding{
			GoogleID: input.UserBinding.GoogleID,
			// FacebookID: input.UserBinding.FacebookID,
			// AppleID:    input.UserBinding.AppleID,
		},
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

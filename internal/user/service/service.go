package usersvc

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	userrepo "github.com/cp25sy5-modjot/main-service/internal/user/repository"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Service interface {
	Create(input *UserCreateInput) (*e.User, error)
	CreateMockUser(input *UserCreateInput, uid string) (*e.User, error)

	GetAll() ([]*e.User, error)
	GetByID(userID string) (*e.User, error)
	GetByGoogleID(google_id string) (*e.User, error)
	// GetByFacebookID(facebook_id string) (*e.User, error)
	// GetByAppleID(apple_id string) (*e.User, error)

	Update(userID string, input *UserUpdateInput) (*e.User, error)
	Delete(userID string) error
	SoftDelete(userID string) error
	TestSoftDelete(userID string) error
	RestoreByGoogleID(googleID string) (*e.User, error)
	RestoreByUserID(userID string) (*e.User, error)
}

type service struct {
	repo        *userrepo.Repository
	asynqClient *asynq.Client
}

func NewService(repo *userrepo.Repository, asynqClient *asynq.Client) *service {
	return &service{repo: repo, asynqClient: asynqClient}
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

func (s *service) GetByID(userID string) (*e.User, error) {
	user, err := s.repo.FindByID(userID)
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

	if exists.Status == e.UserStatusPreActive {
		exists.Status = e.UserStatusActive
	}

	updatedUser, err := s.repo.Update(exists)
	if err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (s *service) Delete(userID string) error {
	return s.repo.HardDelete(userID)
}

func (s *service) SoftDelete(userID string) error {
	// 1. unsubscribe
	if err := s.repo.Unsubscribe(userID); err != nil {
		return err
	}

	// 2. enqueue purge after 30 days
	task, err := tasks.NewPurgeUserTask(
		userID,
		30*24*60*60, // 30 days
	)
	if err != nil {
		return err
	}

	_, err = s.asynqClient.Enqueue(task)
	return err
}

// remove in prod
func (s *service) TestSoftDelete(userID string) error {
	// 1. unsubscribe
	if err := s.repo.Unsubscribe(userID); err != nil {
		return err
	}

	// 2. enqueue purge after 1 mins
	task, err := tasks.NewPurgeUserTask(
		userID,
		1*60, // 1 mins
	)
	if err != nil {
		return err
	}

	_, err = s.asynqClient.Enqueue(task)
	return err
}

func (s *service) RestoreByGoogleID(googleID string) (*e.User, error) {
	// 1. หา user แบบ unscoped (รวม soft deleted)
	user, err := s.repo.FindByGoogleID(googleID)
	if err != nil {
		return nil, err
	}

	// 2. ต้องเป็น inactive เท่านั้น
	if user.Status != e.UserStatusInactive {
		return nil, errors.New("user is not restorable")
	}

	// 3. restore
	if err := s.repo.Restore(user.UserID); err != nil {
		return nil, err
	}

	// 4. return user (ล่าสุด)
	return s.repo.FindByID(user.UserID)
}

func (s *service) RestoreByUserID(userID string) (*e.User, error) {
	// 1. หา user แบบ unscoped (รวม soft deleted)
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 2. ต้องเป็น inactive เท่านั้น
	if user.Status != e.UserStatusInactive {
		return nil, errors.New("user is not restorable")
	}

	// 3. restore
	if err := s.repo.Restore(user.UserID); err != nil {
		return nil, err
	}

	// 4. return user (ล่าสุด)
	return s.repo.FindByID(user.UserID)
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

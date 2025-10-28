package user

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(user *UserInsertReq) error {
	u := &User{
		UserID: uuid.New().String(),
		UserBinding: UserBinding{
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
	return s.repo.Create(u)
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByEmail(email string) (*User, error) {
	return s.repo.FindByEmail(email)
}

func (s *Service) GetByID(user_id string) (*User, error) {
	return s.repo.FindByID(user_id)
}

func (s *Service) GetByGoogleID(google_id string) (*User, error) {
	return s.repo.FindByGoogleID(google_id)
}

func (s *Service) GetByFacebookID(facebook_id string) (*User, error) {
	return s.repo.FindByFacebookID(facebook_id)
}

func (s *Service) GetByAppleID(apple_id string) (*User, error) {
	return s.repo.FindByAppleID(apple_id)
}

func (s *Service) Update(user *User) error {
	return s.repo.Update(user)
}

func (s *Service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

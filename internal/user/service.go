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
	if user.UserID == "" {
		user.UserID = uuid.New().String()
	}
	u := &User{
		UserID:    user.UserID,
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

func (s *Service) Update(user *User) error {
	return s.repo.Update(user)
}

func (s *Service) Delete(user_id string) error {
	return s.repo.Delete(user_id)
}

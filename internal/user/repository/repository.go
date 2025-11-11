package user

import (
	"gorm.io/gorm"
	model "github.com/cp25sy5-modjot/main-service/internal/user/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *model.User) (*model.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAll() ([]*model.User, error) {
	var users []*model.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) FindByID(user_id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("user_id = ?", user_id).First(&user).Error
	return &user, err
}

func (r *Repository) FindByGoogleID(google_id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("google_id = ?", google_id).First(&user).Error
	return &user, err
}

func (r *Repository) FindByFacebookID(facebook_id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("facebook_id = ?", facebook_id).First(&user).Error
	return &user, err
}

func (r *Repository) FindByAppleID(apple_id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("apple_id = ?", apple_id).First(&user).Error
	return &user, err
}

func (r *Repository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(user_id string) error {
	return r.db.Delete(&model.User{}, "user_id = ?", user_id).Error
}

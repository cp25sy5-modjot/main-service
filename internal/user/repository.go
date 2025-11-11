package user

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *User) (*User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAll() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByID(user_id string) (*User, error) {
	var user User
	err := r.db.Where("user_id = ?", user_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByGoogleID(google_id string) (*User, error) {
	var user User
	err := r.db.Where("google_id = ?", google_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByFacebookID(facebook_id string) (*User, error) {
	var user User
	err := r.db.Where("facebook_id = ?", facebook_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByAppleID(apple_id string) (*User, error) {
	var user User
	err := r.db.Where("apple_id = ?", apple_id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Update(user *User) error {
	return r.db.Save(user).Error
}

func (r *Repository) Delete(user_id string) error {
	return r.db.Delete(&User{}, "user_id = ?", user_id).Error
}

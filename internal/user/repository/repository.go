package userrepo

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(user *e.User) (*e.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) FindAll() ([]*e.User, error) {
	var users []*e.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *Repository) FindByID(userID string) (*e.User, error) {
	var user e.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (r *Repository) FindByGoogleID(googleID string) (*e.User, error) {
	var user e.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	return &user, err
}

// func (r *Repository) FindByFacebookID(facebook_id string) (*e.User, error) {
// 	var user e.User
// 	err := r.db.Where("facebook_id = ?", facebook_id).First(&user).Error
// 	return &user, err
// }

// func (r *Repository) FindByAppleID(apple_id string) (*e.User, error) {
// 	var user e.User
// 	err := r.db.Where("apple_id = ?", apple_id).First(&user).Error
// 	return &user, err
// }

func (r *Repository) Update(user *e.User) (*e.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) HardDelete(userID string) error {
	return r.db.Unscoped().
		Delete(&e.User{}, "user_id = ?", userID).
		Error
}

func (r *Repository) Unsubscribe(id string) error {
	now := time.Now()

	tx := r.db.Model(&e.User{}).
		Where("user_id = ?", id).
		Where("status != ?", e.StatusInactive).
		Updates(map[string]interface{}{
			"status":          e.StatusInactive,
			"unsubscribed_at": &now,
		})

	if tx.RowsAffected == 0 {
		return nil
	}
	return tx.Error
}

func (r *Repository) PurgeExpiredUnsubscribed(days int) error {
	return r.db.
		Where("status = ?", e.StatusInactive).
		Where("unsubscribed_at < ?", time.Now().AddDate(0, 0, -days)).
		Unscoped().
		Delete(&e.User{}).
		Error
}

func (r *Repository) Restore(id string) error {
	return r.db.Model(&e.User{}).
		Where("user_id = ?", id).
		Updates(map[string]interface{}{
			"status":          e.StatusActive,
			"unsubscribed_at": nil,
		}).Error
}

func (r *Repository) RestoreAndReturn(userID string) (*e.User, error) {
	var user e.User

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// lock row
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		if user.Status != e.StatusInactive {
			return errors.New("not restorable")
		}

		if err := tx.Model(&e.User{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"status":          e.StatusActive,
				"unsubscribed_at": nil,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	return &user, err
}

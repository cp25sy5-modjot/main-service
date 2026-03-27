package userrepo

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(user *e.User) (*e.User, error)

	FindAll() ([]*e.User, error)

	FindByID(userID string) (*e.User, error)

	FindByGoogleID(googleID string) (*e.User, error)

	Update(user *e.User) (*e.User, error)

	HardDelete(userID string) error

	Unsubscribe(id string) error

	PurgeExpiredUnsubscribed(days int) error

	Restore(id string) error

	RestoreAndReturn(userID string) (*e.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(user *e.User) (*e.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) FindAll() ([]*e.User, error) {
	var users []*e.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *repository) FindByID(userID string) (*e.User, error) {
	var user e.User
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}

func (r *repository) FindByGoogleID(googleID string) (*e.User, error) {
	var user e.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	return &user, err
}

func (r *repository) Update(user *e.User) (*e.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) HardDelete(userID string) error {
	return r.db.Unscoped().
		Delete(&e.User{}, "user_id = ?", userID).
		Error
}

func (r *repository) Unsubscribe(id string) error {
	now := time.Now()

	tx := r.db.Model(&e.User{}).
		Where("user_id = ?", id).
		Where("status != ?", e.UserStatusInactive).
		Updates(map[string]interface{}{
			"status":          e.UserStatusInactive,
			"unsubscribed_at": &now,
		})

	if tx.RowsAffected == 0 {
		return nil
	}
	return tx.Error
}

func (r *repository) PurgeExpiredUnsubscribed(days int) error {
	return r.db.
		Where("status = ?", e.UserStatusInactive).
		Where("unsubscribed_at < ?", time.Now().AddDate(0, 0, -days)).
		Unscoped().
		Delete(&e.User{}).
		Error
}

func (r *repository) Restore(id string) error {
	return r.db.Model(&e.User{}).
		Where("user_id = ?", id).
		Updates(map[string]interface{}{
			"status":          e.UserStatusActive,
			"unsubscribed_at": nil,
		}).Error
}

func (r *repository) RestoreAndReturn(userID string) (*e.User, error) {
	var user e.User

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// lock row
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		if user.Status != e.UserStatusInactive {
			return errors.New("not restorable")
		}

		if err := tx.Model(&e.User{}).
			Where("user_id = ?", userID).
			Updates(map[string]interface{}{
				"status":          e.UserStatusActive,
				"unsubscribed_at": nil,
			}).Error; err != nil {
			return err
		}

		return nil
	})

	return &user, err
}

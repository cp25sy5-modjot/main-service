package pushrepo

import (
	"context"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Save(ctx context.Context, userID, token, platform string) error

	FindByUserID(ctx context.Context, userID string) ([]string, error)

	DeleteByToken(ctx context.Context, token string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Save(ctx context.Context, userID, token, platform string) error {
	pt := e.PushToken{
		ID:       uuid.New().String(),
		UserID:   userID,
		Token:    token,
		Platform: platform,
	}

	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "token"}},
			DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform", "updated_at"}),
		}).
		Create(&pt).Error
}

func (r *repository) FindByUserID(ctx context.Context, userID string) ([]string, error) {
	var tokens []string

	err := r.db.WithContext(ctx).
		Model(&e.PushToken{}).
		Where("user_id = ?", userID).
		Pluck("token", &tokens).Error

	return tokens, err
}

func (r *repository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).
		Where("token = ?", token).
		Delete(&e.PushToken{}).Error
}

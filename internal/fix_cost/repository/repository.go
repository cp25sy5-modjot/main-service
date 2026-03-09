package fixcostrepo

import (
	"context"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) WithTx(tx *gorm.DB) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) Create(ctx context.Context, fc *e.FixCost) error {
	return r.db.WithContext(ctx).Create(fc).Error
}

func (r *Repository) Update(ctx context.Context, fc *e.FixCost) error {
	return r.db.WithContext(ctx).Save(fc).Error
}

func (r *Repository) Delete(ctx context.Context, id string, userID string) error {
	return r.db.WithContext(ctx).Delete(&e.FixCost{}, "fix_cost_id = ? AND user_id = ?", id, userID).Error
}

func (r *Repository) FindAllActive(ctx context.Context) ([]e.FixCost, error) {
	var fixCosts []e.FixCost

	err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Find(&fixCosts).Error

	if err != nil {
		return nil, err
	}

	return fixCosts, nil
}

func (r *Repository) FindByID(ctx context.Context, id string, userID string) (*e.FixCost, error) {
	var fc e.FixCost

	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("fix_cost_id = ? AND user_id = ?", id, userID).
		Order("created_at DESC").
		First(&fc).Error

	if err != nil {
		return nil, err
	}

	return &fc, nil
}

func (r *Repository) FindAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error) {
	var fixCosts []*e.FixCost

	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&fixCosts).Error

	if err != nil {
		return nil, err
	}

	return fixCosts, nil
}

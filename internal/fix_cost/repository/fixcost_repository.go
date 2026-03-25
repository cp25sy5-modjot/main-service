package fixcostrepo

import (
	"context"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, fc *e.FixCost) error
	Update(ctx context.Context, fc *e.FixCost) error
	Delete(ctx context.Context, id string, userID string) error

	FindAllActive(ctx context.Context) ([]e.FixCost, error)
	FindByID(ctx context.Context, id string, userID string) (*e.FixCost, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error)
	FindDueFixCosts(ctx context.Context) ([]*e.FixCost, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) WithTx(tx *gorm.DB) *repository {
	return &repository{db: tx}
}

func (r *repository) Create(ctx context.Context, fc *e.FixCost) error {
	return r.db.WithContext(ctx).Create(fc).Error
}

func (r *repository) Update(ctx context.Context, fc *e.FixCost) error {
	return r.db.WithContext(ctx).Save(fc).Error
}

func (r *repository) Delete(ctx context.Context, id string, userID string) error {
	return r.db.WithContext(ctx).Delete(&e.FixCost{}, "fix_cost_id = ? AND user_id = ?", id, userID).Error
}

func (r *repository) FindAllActive(ctx context.Context) ([]e.FixCost, error) {
	var fixCosts []e.FixCost

	err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Find(&fixCosts).Error

	if err != nil {
		return nil, err
	}

	return fixCosts, nil
}

func (r *repository) FindByID(ctx context.Context, id string, userID string) (*e.FixCost, error) {
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

func (r *repository) FindAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error) {
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

func (r *repository) FindDueFixCosts(ctx context.Context) ([]*e.FixCost, error) {
	var fcs []*e.FixCost

	err := r.db.WithContext(ctx).
		Where("status = ? AND next_run_date <= ?", "active", time.Now().UTC()).
		Order("next_run_date ASC").
		Limit(100).
		Find(&fcs).Error

	return fcs, err
}

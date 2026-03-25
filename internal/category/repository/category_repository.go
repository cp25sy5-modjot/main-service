package categoryrepo

import (
	"context"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"gorm.io/gorm"
)

type Repository interface {
	Create(category *e.Category) (*e.Category, error)
	FindAllByUserID(userID string) ([]e.Category, error)
	FindByID(params *m.CategorySearchParams) (*e.Category, error)
	FindAllByUserIDWithTransactionsFiltered(userID string, start, end time.Time) ([]e.Category, error)
	FindByIDWithTransactionsFiltered(params *m.CategorySearchParams, start, end time.Time) (*e.Category, error)
	Update(category *e.Category) (*e.Category, error)
	Delete(params *m.CategorySearchParams) error
	FindByIDs(ctx context.Context, userID string, ids []string) (map[string]e.Category, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) Create(category *e.Category) (*e.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *repository) FindAllByUserID(userID string) ([]e.Category, error) {
	var categories []e.Category
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *repository) FindByID(params *m.CategorySearchParams) (*e.Category, error) {
	var category e.Category
	err := r.db.
		First(&category,
			"category_id = ? AND user_id = ?",
			params.CategoryID,
			params.UserID).Error
	return &category, err
}

func (r *repository) FindAllByUserIDWithTransactionsFiltered(
	userID string,
	start, end time.Time,
) ([]e.Category, error) {
	var categories []e.Category

	err := r.db.
		Preload("TransactionItems", func(tx *gorm.DB) *gorm.DB {
			return tx.
				Joins("JOIN transactions t ON t.transaction_id = transaction_items.transaction_id").
				Where("t.user_id = ? AND t.date >= ? AND t.date < ?", userID, start, end)
		}).
		Preload("TransactionItems.Transaction").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error

	return categories, err
}

func (r *repository) FindByIDWithTransactionsFiltered(
	params *m.CategorySearchParams,
	start, end time.Time,
) (*e.Category, error) {
	var category e.Category

	err := r.db.
		Preload("TransactionItems", func(tx *gorm.DB) *gorm.DB {
			return tx.
				Joins("JOIN transactions t ON t.transaction_id = transaction_items.transaction_id").
				Where(
					"t.user_id = ? AND t.date >= ? AND t.date < ?",
					params.UserID,
					start,
					end,
				)
		}).
		Preload("TransactionItems.Transaction").
		First(
			&category,
			"category_id = ? AND user_id = ?",
			params.CategoryID,
			params.UserID,
		).Error

	return &category, err
}

func (r *repository) Update(category *e.Category) (*e.Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *repository) Delete(params *m.CategorySearchParams) error {
	return r.db.Delete(&e.Category{}, "category_id = ? AND user_id = ?", params.CategoryID, params.UserID).Error
}

func (r *repository) FindByIDs(
	ctx context.Context,
	userID string,
	ids []string,
) (map[string]e.Category, error) {

	result := make(map[string]e.Category)

	if len(ids) == 0 {
		return result, nil
	}

	var categories []e.Category

	if err := r.db.
		WithContext(ctx).
		Where("user_id = ? AND category_id IN ?", userID, ids).
		Find(&categories).Error; err != nil {
		return nil, err
	}

	for _, c := range categories {
		result[c.CategoryID] = c
	}

	return result, nil
}

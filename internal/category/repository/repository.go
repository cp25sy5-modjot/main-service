package categoryrepo

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(category *e.Category) (*e.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *Repository) FindAllByUserID(userID string) ([]e.Category, error) {
	var categories []e.Category
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *Repository) FindByID(params *m.CategorySearchParams) (*e.Category, error) {
	var category e.Category
	err := r.db.
		First(&category,
			"category_id = ? AND user_id = ?",
			params.CategoryID,
			params.UserID).Error
	return &category, err
}

func (r *Repository) FindAllByUserIDWithTransactions(userID string) ([]e.Category, error) {
	var categories []e.Category
	err := r.db.
		Preload("Transactions").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *Repository) FindByIDWithTransactions(params *m.CategorySearchParams) (*e.Category, error) {
	var category e.Category
	err := r.db.
		Preload("Transactions").
		First(&category,
			"category_id = ? AND user_id = ?",
			params.CategoryID,
			params.UserID).Error
	return &category, err
}

func (r *Repository) Update(category *e.Category) (*e.Category, error) {
	if err := r.db.Save(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *Repository) Delete(params *m.CategorySearchParams) error {
	return r.db.Delete(&e.Category{}, "category_id = ? AND user_id = ?", params.CategoryID, params.UserID).Error
}

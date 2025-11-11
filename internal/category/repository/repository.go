package category

import (
	model "github.com/cp25sy5-modjot/main-service/internal/category/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(category *model.Category) (*model.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *Repository) FindAllByUserID(userID string) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&categories).Error
	return categories, err
}

func (r *Repository) FindByID(params *model.CategorySearchParams) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, "category_id = ? AND user_id = ?", params.CategoryID, params.UserID).Error
	return &category, err
}

func (r *Repository) Update(category *model.Category) error {

	return r.db.Save(category).Error
}

func (r *Repository) Delete(params *model.CategorySearchParams) error {
	return r.db.Delete(&model.Category{}, "category_id = ? AND user_id = ?", params.CategoryID, params.UserID).Error
}

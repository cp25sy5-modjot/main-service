package categoryrepo

import (
	"time"

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

func (r *Repository) FindAllByUserIDWithTransactionsFiltered(
	userID string,
	start, end time.Time,
) ([]e.Category, error) {
	var categories []e.Category

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error

	return categories, err
}

func (r *Repository) FindByIDWithTransactionsFiltered(
	params *m.CategorySearchParams,
	start, end time.Time,
) (*e.Category, error) {
	var category e.Category

	err := r.db.
		First(
			&category,
			"category_id = ? AND user_id = ?",
			params.CategoryID,
			params.UserID,
		).Error

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

func (r *Repository) GetCategoriesAndTransactions(
	userId string,
	start, end time.Time,
) ([]m.CategoryRes, error) {

	var list []m.CategoryRes

	err := r.db.Raw(`
		SELECT 
			c.category_id,
			c.category_name,
			c.budget,
			c.color_code,
			c.created_at,
			COALESCE(SUM(t.price), 0) AS budget_usage
		FROM categories c
		LEFT JOIN transactions t
			ON t.category_id = c.category_id
			AND t.user_id = c.user_id
			AND t.date >= ?
			AND t.date < ?
		WHERE c.user_id = ?
		GROUP BY c.category_id, c.category_name, c.color_code, c.budget
		ORDER BY budget_usage DESC
	`, start, end, userId).Scan(&list).Error

	return list, err
}

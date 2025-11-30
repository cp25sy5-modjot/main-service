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
	filter *m.TransactionFilter,
) ([]e.Category, error) {
	var categories []e.Category

	startOfMonth, endOfMonth := getMonthRange(filter.Date)

	err := r.db.
		// preload เฉพาะ transactions ที่อยู่ในช่วงวันที่ที่ต้องการ
		Preload("Transactions", "date >= ? AND date < ?", startOfMonth, endOfMonth).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&categories).Error

	return categories, err
}

func (r *Repository) FindByIDWithTransactionsFiltered(
	params *m.CategorySearchParams,
	filter *m.TransactionFilter,
) (*e.Category, error) {
	var category e.Category

	startOfMonth, endOfMonth := getMonthRange(filter.Date)

	err := r.db.
		Preload("Transactions", "date >= ? AND date < ?", startOfMonth, endOfMonth).
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

func getMonthRange(t *time.Time) (time.Time, time.Time) {
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return startOfMonth, endOfMonth
}

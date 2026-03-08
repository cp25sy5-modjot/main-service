package summaryrepo

import (
	"context"
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) ExpenseSummary(
	ctx context.Context,
	userID string,
	format string,
	start string,
	end string,
) ([]m.ExpenseSummary, error) {

	var result []m.ExpenseSummary

	err := r.db.WithContext(ctx).
		Table("transaction_items ti").
		Select(`
			DATE_FORMAT(tr.date, ?) AS label,
			COALESCE(SUM(ti.price),0) AS total
		`, format).
		Joins("JOIN transactions tr ON tr.transaction_id = ti.transaction_id").
		Where("tr.user_id = ?", userID).
		Where("tr.date >= ? AND tr.date < ?", start, end).
		Group("label").
		Order("label").
		Scan(&result).Error

	return result, err
}

func (r *Repository) CategorySummary(
	ctx context.Context,
	userID string,
	start time.Time,
	end time.Time,
) ([]m.CategorySummary, error) {

	var result []m.CategorySummary

	err := r.db.WithContext(ctx).
		Table("transaction_items ti").
		Select(`
			c.category_id,
			c.icon AS category_icon,
			c.name AS category_name,
			c.color AS category_color,
			COALESCE(SUM(ti.price),0) AS total
		`).
		Joins("JOIN transactions tr ON tr.transaction_id = ti.transaction_id").
		Joins("JOIN categories c ON c.category_id = ti.category_id").
		Where("tr.user_id = ?", userID).
		Where("tr.date >= ? AND tr.date < ?", start, end).
		Group("c.category_id, c.icon, c.name, c.color").
		Order("total DESC").
		Scan(&result).Error

	return result, err
}

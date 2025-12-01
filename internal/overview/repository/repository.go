package overviewrepo

import (
	"time"

	"gorm.io/gorm"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

// internal/overview/repository.go (cont.)

// --- Last N transactions (price already = price * quantity) ---
func (r *Repository) GetLastTransactions(userID string, start, end time.Time, limit int) ([]m.LastTransaction, error) {
	var list []m.LastTransaction

	err := r.db.
		Table("transactions t").
		Select(`
			t.transaction_id,
			t.item_id,
			t.title,
			t.price,
			t.date,
			t.type,
			t.category_id,
			COALESCE(c.category_name, '') AS category_name,
			COALESCE(c.color_code, '')   AS category_color_code
		`).
		Joins("LEFT JOIN categories c ON c.category_id = t.category_id AND c.user_id = t.user_id").
		Where("t.user_id = ? AND t.date >= ? AND t.date < ?", userID, start, end).
		Order("t.date DESC").
		Limit(limit).
		Scan(&list).Error

	return list, err
}

// --- Top categories by spending in month (budget_usage = % of budget) ---
func (r *Repository) GetTopCategoriesBySpending(
	userID string,
	start, end time.Time,
	limit int,
) ([]m.TopCategoryUsage, error) {

	var list []m.TopCategoryUsage

	err := r.db.Raw(`
		SELECT 
			c.category_id,
			c.category_name,
			c.color_code,
			c.budget,
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
		LIMIT ?
	`, start, end, userID, limit).Scan(&list).Error

	return list, err
}

func (r *Repository) GetMonthTotal(userID string, start, end time.Time) (float64, error) {
	var total float64

	err := r.db.
		Table("transactions").
		Select("COALESCE(SUM(price), 0)").
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Scan(&total).Error

	return total, err
}

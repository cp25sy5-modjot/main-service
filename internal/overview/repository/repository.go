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
func (r *Repository) GetLastTransactions(userID string, limit int) ([]m.LastTransaction, error) {
	var list []m.LastTransaction

	err := r.db.
		Table("transaction_items ti").
		Select(`
			ti.transaction_id,
			ti.item_id,
			ti.title,
			ti.price,
			tr.date,
			tr.type,
			ti.category_id,
			COALESCE(c.category_name, '') AS category_name,
			COALESCE(c.color_code, '')   AS category_color_code
		`).
		Joins("JOIN transactions tr ON tr.transaction_id = ti.transaction_id").
		Joins("LEFT JOIN categories c ON c.category_id = ti.category_id AND c.user_id = tr.user_id").
		Where("tr.user_id = ?", userID).
		Order("tr.date DESC").
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
			COALESCE(SUM(ti.price), 0) AS budget_usage
		FROM categories c
		LEFT JOIN (
			SELECT ti.*
			FROM transaction_items ti
			JOIN transactions tr
				ON tr.transaction_id = ti.transaction_id
			WHERE tr.user_id = ?
			  AND tr.date >= ?
			  AND tr.date < ?
		) ti
			ON ti.category_id = c.category_id
		WHERE c.user_id = ?
		GROUP BY c.category_id, c.category_name, c.color_code, c.budget
		ORDER BY budget_usage DESC
		LIMIT ?
	`, userID, start, end, userID, limit).Scan(&list).Error

	return list, err
}


func (r *Repository) GetMonthTotal(userID string, start, end time.Time) (float64, error) {
	var total float64

	err := r.db.
		Table("transaction_items ti").
		Select("COALESCE(SUM(ti.price), 0)").
		Joins("JOIN transactions tr ON tr.transaction_id = ti.transaction_id").
		Where("tr.user_id = ? AND tr.date >= ? AND tr.date < ?", userID, start, end).
		Scan(&total).Error

	return total, err
}


package transaction

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

func (r *Repository) Create(transaction *e.Transaction) (*e.Transaction, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *Repository) FindAllByUserID(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Preload("Category"). // load related Category
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindAllByUserIDAndFiltered(userID string, filter *m.TransactionFilter) ([]e.Transaction, error) {
	var transactions []e.Transaction

	t := filter.Date

	// Calculate the start of the target month (e.g., 2025-11-01 00:00:00)
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	// Calculate the start of the next month (e.g., 2025-12-01 00:00:00)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	// Filter by user_id AND date >= startOfMonth AND date < endOfMonth
	err := r.db.
		Preload("Category").
		Where("user_id = ? AND date >= ? AND date < ?",
			userID,
			startOfMonth,
			endOfMonth,
		).
		Order("date DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *Repository) FindByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	var transaction e.Transaction
	err := r.db.
		Preload("Category").
		First(&transaction,
			"transaction_id = ? AND item_id = ? AND user_id = ?",
			params.TransactionID,
			params.ItemID,
			params.UserID).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *e.Transaction) (*e.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *Repository) Delete(params *m.TransactionSearchParams) error {
	return r.db.Delete(&e.Transaction{}, "transaction_id = ? AND item_id = ?", params.TransactionID, params.ItemID).Error
}

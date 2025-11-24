package transaction

import (
	"log"
	"time"

	model "github.com/cp25sy5-modjot/main-service/internal/transaction/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(transaction *model.Transaction) (*model.Transaction, error) {
	if err := r.db.Create(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *Repository) FindAllByUserID(userID string) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("user_id = ?", userID).Order("date DESC").Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindAllByUserIDAndFiltered(userID string, filter *model.TransactionFilter) ([]model.Transaction, error) {
	var transactions []model.Transaction

	t := filter.Date

	// Calculate the start of the target month (e.g., 2025-11-01 00:00:00)
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

	// Calculate the start of the next month (e.g., 2025-12-01 00:00:00)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	log.Printf("Filtering transactions from %s to %s", startOfMonth.String(), endOfMonth.String())
	// Filter by user_id AND date >= startOfMonth AND date < endOfMonth
	err := r.db.Where("user_id = ? AND date >= ? AND date < ?",
		userID,
		startOfMonth,
		endOfMonth,
	).Order("date DESC").Find(&transactions).Error

	return transactions, err
}

func (r *Repository) FindByID(params *model.TransactionSearchParams) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.First(&transaction, "transaction_id = ? AND product_id = ? AND user_id = ?", params.TransactionID, params.ItemID, params.UserID).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *model.Transaction) (*model.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *Repository) Delete(params *model.TransactionSearchParams) error {
	return r.db.Delete(&model.Transaction{}, "transaction_id = ? AND product_id = ?", params.TransactionID, params.ItemID).Error
}

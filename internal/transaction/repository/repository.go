package transactionrepo

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

func (r *Repository) findAllUncategorizedByUserID(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Preload("Category"). // load related Category
		Where("user_id = ? AND category_id IS NULL", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
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

    // Start of current month
    startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())

    var startRange, endRange time.Time

    if filter.PreviousMonth {
        // Previous month
        startRange = startOfMonth.AddDate(0, -1, 0)
        endRange   = startOfMonth
    } else {
        // Current month
        startRange = startOfMonth
        endRange   = startOfMonth.AddDate(0, 1, 0)
    }

    err := r.db.
        Preload("Category").
        Where("user_id = ? AND date >= ? AND date < ?", userID, startRange, endRange).
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

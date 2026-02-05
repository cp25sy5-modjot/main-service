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

func (r *Repository) WithTx(tx *gorm.DB) *Repository {
	return &Repository{db: tx}
}

func (r *Repository) Create(transaction *e.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *Repository) FindAllByUserID(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindAllByUserIDWithRelations(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Preload("Items").
		Preload("Items.Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindAllByUserIDAndFiltered(userID string, start, end time.Time) ([]e.Transaction, error) {
	var transactions []e.Transaction

	err := r.db.
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Order("date DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *Repository) FindAllByUserIDWithRelationsAndFiltered(
	userID string,
	start, end time.Time,
	categoryIDs []string,
) ([]e.Transaction, error) {

	var transactions []e.Transaction

	query := r.db.
		Model(&e.Transaction{}).
		Preload("Items").
		Preload("Items.Category").
		Where(
			"transactions.user_id = ? AND transactions.date >= ? AND transactions.date < ?",
			userID,
			start,
			end,
		)

	// 👉 filter ด้วย multiple category จาก item
	if len(categoryIDs) > 0 {
		query = query.
			Joins("JOIN transaction_items ON transaction_items.transaction_id = transactions.transaction_id").
			Where("transaction_items.category_id IN ?", categoryIDs).
			Distinct()
	}

	err := query.
		Order("transactions.date DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *Repository) CountItemsByUserAndDateRange(
	userID string,
	start, end time.Time,
) (int, error) {

	var count int64

	err := r.db.
		Table("transaction_items").
		Joins("JOIN transactions ON transactions.transaction_id = transaction_items.transaction_id").
		Where(
			"transactions.user_id = ? AND transactions.date >= ? AND transactions.date < ?",
			userID,
			start,
			end,
		).
		Count(&count).Error

	return int(count), err
}

func (r *Repository) FindByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	var transaction e.Transaction
	err := r.db.
		First(&transaction,
			"transaction_id = ? AND user_id = ?",
			params.TransactionID,
			params.UserID).Error
	return &transaction, err
}

func (r *Repository) FindByIDWithRelations(
	params *m.TransactionSearchParams,
) (*e.Transaction, error) {

	var transaction e.Transaction
	err := r.db.
		Preload("Items").
		Preload("Items.Category").
		First(&transaction,
			"transaction_id = ? AND user_id = ?",
			params.TransactionID,
			params.UserID,
		).Error

	return &transaction, err
}

func (r *Repository) Update(transaction *e.Transaction) (*e.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	return r.FindByIDWithRelations(&m.TransactionSearchParams{
		TransactionID: transaction.TransactionID,
		UserID:        transaction.UserID,
	})
}

func (r *Repository) Delete(params *m.TransactionSearchParams) error {
	return r.db.Delete(&e.Transaction{}, "transaction_id = ? AND user_id = ?", params.TransactionID, params.UserID).Error
}

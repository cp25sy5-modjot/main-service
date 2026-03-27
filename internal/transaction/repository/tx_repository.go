package transactionrepo

import (
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"gorm.io/gorm"
)

type Repository interface {
	WithTx(tx *gorm.DB) Repository
	Create(transaction *e.Transaction) error
	FindAllByUserID(userID string) ([]e.Transaction, error)
	FindAllByUserIDWithRelations(userID string) ([]e.Transaction, error)
	FindAllByUserIDAndFiltered(userID string, start, end time.Time) ([]e.Transaction, error)
	FindAllByUserIDWithRelationsAndFiltered(userID string, start, end time.Time, categoryIDs []string) ([]e.Transaction, error)
	CountItemsByUserAndDateRange(userID string, start, end time.Time) (int, error)
	FindByID(params *m.TransactionSearchParams) (*e.Transaction, error)
	FindByIDWithRelations(params *m.TransactionSearchParams) (*e.Transaction, error)
	Update(transaction *e.Transaction) (*e.Transaction, error)
	UpdateFieldsTx(tx *gorm.DB, txID string, updates map[string]interface{}) error
	Delete(params *m.TransactionSearchParams) error
	FindByFixCostIDAndRunDate(params *m.TransactionFixCostSearchParams) (*e.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) WithTx(tx *gorm.DB) Repository {
	return &repository{db: tx}
}

func (r *repository) Create(transaction *e.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *repository) FindAllByUserID(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *repository) FindAllByUserIDWithRelations(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
		Preload("Items").
		Preload("Items.Category").
		Where("user_id = ?", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *repository) FindAllByUserIDAndFiltered(userID string, start, end time.Time) ([]e.Transaction, error) {
	var transactions []e.Transaction

	err := r.db.
		Where("user_id = ? AND date >= ? AND date < ?", userID, start, end).
		Order("date DESC").
		Find(&transactions).Error

	return transactions, err
}

func (r *repository) FindAllByUserIDWithRelationsAndFiltered(
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

func (r *repository) CountItemsByUserAndDateRange(
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

func (r *repository) FindByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	var transaction e.Transaction
	err := r.db.
		First(&transaction,
			"transaction_id = ? AND user_id = ?",
			params.TransactionID,
			params.UserID).Error
	return &transaction, err
}

func (r *repository) FindByIDWithRelations(
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

func (r *repository) Update(transaction *e.Transaction) (*e.Transaction, error) {
	if err := r.db.Save(transaction).Error; err != nil {
		return nil, err
	}

	return r.FindByIDWithRelations(&m.TransactionSearchParams{
		TransactionID: transaction.TransactionID,
		UserID:        transaction.UserID,
	})
}

func (r *repository) UpdateFieldsTx(
	tx *gorm.DB,
	txID string,
	updates map[string]interface{},
) error {
	return tx.
		Model(&e.Transaction{}).
		Where("transaction_id = ?", txID).
		Updates(updates).
		Error
}

func (r *repository) Delete(params *m.TransactionSearchParams) error {
	return r.db.Delete(&e.Transaction{}, "transaction_id = ? AND user_id = ?", params.TransactionID, params.UserID).Error
}

func (r *repository) FindByFixCostIDAndRunDate(
	params *m.TransactionFixCostSearchParams,
) (*e.Transaction, error) {

	var tx e.Transaction

	runDate := params.RunDate.UTC().Truncate(24 * time.Hour)

	err := r.db.
		Preload("Items").
		Where("fix_cost_id = ? AND run_date = ? AND user_id = ?",
			params.FixCostID,
			runDate,
			params.UserID,
		).
		First(&tx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil 
		}
		return nil, err
	}

	return &tx, nil
}

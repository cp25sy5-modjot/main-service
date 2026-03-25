package txirepo

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"gorm.io/gorm"
)

type Repository interface {
	WithTx(fn func(tx *gorm.DB) error) error
	Create(items *e.TransactionItem) error
	FindAllByUserID(params *m.TransactionItemSearchParams) ([]e.TransactionItem, error)
	FindByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error)
	FindByIDWithCategory(params *m.TransactionItemSearchParams) (*e.TransactionItem, error)
	Update(items *e.TransactionItem) (*e.TransactionItem, error)
	DeleteItem(params *m.TransactionItemSearchParams) error
	CreateManyTx(tx *gorm.DB, items []e.TransactionItem) error
	DeleteByTransactionIDTx(tx *gorm.DB, transactionID string) error
	CountItemsByTransactionID(transactionID string) (int64, error)
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) WithTx(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (r *repository) Create(items *e.TransactionItem) error {
	if err := r.db.
		Preload("Category").
		Create(items).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAllByUserID(params *m.TransactionItemSearchParams) ([]e.TransactionItem, error) {
	var items []e.TransactionItem
	err := r.db.
		Where("user_id = ?", params.UserID).
		Find(&items).Error
	return items, err
}

func (r *repository) FindByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error) {
	var items e.TransactionItem
	err := r.db.
		First(&items,
			"transaction_id = ? AND item_id = ?",
			params.TransactionID,
			params.ItemID).Error
	return &items, err
}

func (r *repository) FindByIDWithCategory(params *m.TransactionItemSearchParams) (*e.TransactionItem, error) {
	var items e.TransactionItem
	err := r.db.
		Preload("Category").
		First(&items,
			"transaction_id = ? AND item_id = ?",
			params.TransactionID,
			params.ItemID).Error
	return &items, err
}

func (r *repository) Update(items *e.TransactionItem) (*e.TransactionItem, error) {
	if err := r.db.
		Model(&e.TransactionItem{}).
		Where("transaction_id = ? AND item_id = ?", items.TransactionID, items.ItemID).
		Updates(items).
		Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Preload("Category").
		First(items,
			"transaction_id = ? AND item_id = ?",
			items.TransactionID,
			items.ItemID,
		).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *repository) DeleteItem(params *m.TransactionItemSearchParams) error {
	return r.db.
		Delete(
			&e.TransactionItem{},
			"transaction_id = ? AND item_id = ?",
			params.TransactionID,
			params.ItemID,
		).
		Error
}

func (r *repository) CreateManyTx(
	tx *gorm.DB,
	items []e.TransactionItem,
) error {
	if len(items) == 0 {
		return nil
	}
	return tx.Create(&items).Error
}

func (r *repository) DeleteByTransactionIDTx(
	tx *gorm.DB,
	transactionID string,
) error {
	return tx.
		Where("transaction_id = ?", transactionID).
		Delete(&e.TransactionItem{}).
		Error
}

func (r *repository) CountItemsByTransactionID(transactionID string) (int64, error) {
	var count int64

	err := r.db.
		Model(&e.TransactionItem{}).
		Where("transaction_id = ?", transactionID).
		Count(&count).
		Error

	return count, err
}

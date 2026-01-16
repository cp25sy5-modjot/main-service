package transactionrepo

import (
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

func (r *Repository) Create(items *e.TransactionItem) error {
	if err := r.db.
		Preload("Category").
		Create(items).
		Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) CreateMany(items []e.TransactionItem) error {
	if len(items) == 0 {
		return nil
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&items).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *Repository) FindAllByUserID(params *m.TransactionItemSearchParams) ([]e.TransactionItem, error) {
	var items []e.TransactionItem
	err := r.db.
		Where("user_id = ?", params.UserID).
		Find(&items).Error
	return items, err
}

func (r *Repository) FindByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error) {
	var items e.TransactionItem
	err := r.db.
		First(&items,
			"transaction_id = ? AND item_id = ?",
			params.TransactionID,
			params.ItemID).Error
	return &items, err
}

func (r *Repository) FindByIDWithCategory(params *m.TransactionItemSearchParams) (*e.TransactionItem, error) {
	var items e.TransactionItem
	err := r.db.
		Preload("Category").
		First(&items,
			"transaction_id = ? AND item_id = ?",
			params.TransactionID,
			params.ItemID).Error
	return &items, err
}

func (r *Repository) Update(items *e.TransactionItem) (*e.TransactionItem, error) {
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

func (r *Repository) Delete(params *m.TransactionItemSearchParams) error {
	return r.db.Delete(&e.TransactionItem{}, "transaction_id = ? AND item_id = ?", params.TransactionID, params.ItemID).Error
}

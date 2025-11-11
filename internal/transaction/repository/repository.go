package transaction

import (
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

func (r *Repository) FindByID(params *model.TransactionSearchParams) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.First(&transaction, "transaction_id = ? AND product_id = ? AND user_id = ?", params.TransactionID, params.ItemID, params.UserID).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *model.Transaction) error {

	return r.db.Save(transaction).Error
}

func (r *Repository) Delete(params *model.TransactionSearchParams) error {
	return r.db.Delete(&model.Transaction{}, "transaction_id = ? AND product_id = ?", params.TransactionID, params.ItemID).Error
}

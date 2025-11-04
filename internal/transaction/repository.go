package transaction

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(transaction *Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *Repository) FindAllByUserID(userID string) ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindByID(params *SearchParams) (*Transaction, error) {
	var transaction Transaction
	err := r.db.First(&transaction, "transaction_id = ? AND product_id = ? AND user_id = ?", params.TransactionID, params.ItemID, params.UserID).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *Transaction) error {

	return r.db.Save(transaction).Error
}

func (r *Repository) Delete(params *SearchParams) error {
	return r.db.Delete(&Transaction{}, "transaction_id = ? AND product_id = ?", params.TransactionID, params.ItemID).Error
}

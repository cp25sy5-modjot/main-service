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

func (r *Repository) FindAll() ([]Transaction, error) {
	var transactions []Transaction
	err := r.db.Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindByID(tx_id string, prod_id string) (*Transaction, error) {
	var transaction Transaction
	err := r.db.First(&transaction, "transaction_id = ? AND product_id = ?", tx_id, prod_id).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *Repository) Delete(tx_id string, prod_id string) error {
	return r.db.Delete(&Transaction{}, "transaction_id = ? AND product_id = ?", tx_id, prod_id).Error
}

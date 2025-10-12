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

func (r *Repository) FindByID(id uint) (*Transaction, error) {
	var transaction Transaction
	err := r.db.First(&transaction, id).Error
	return &transaction, err
}

func (r *Repository) Update(transaction *Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *Repository) Delete(id uint) error {
	return r.db.Delete(&Transaction{}, id).Error
}

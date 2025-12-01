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
		Where("user_id = ? AND category_id IS NULL", userID).
		Order("date DESC").
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) FindAllByUserID(userID string) ([]e.Transaction, error) {
	var transactions []e.Transaction
	err := r.db.
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

func (r *Repository) FindByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	var transaction e.Transaction
	err := r.db.
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

func (r *Repository) GetTransactionsWithCategory(
	userId string,
	start, end time.Time,
) ([]m.TransactionRes, error) {
	var list []m.TransactionRes

	err := r.db.Raw(`
		SELECT 
			t.transaction_id,
			t.item_id,
			t.user_id,
			t.title,
			t.price,
			t.quantity,
			t.date,
			t.type,
			t.category_id,
			c.category_name,
			c.color_code,
			c.budget
		FROM transactions t
		LEFT JOIN categories c 
			ON c.category_id = t.category_id
		   AND c.user_id = t.user_id
		WHERE t.user_id = ? 
		  AND t.date >= ? AND t.date < ?
		ORDER BY t.date DESC, t.transaction_id, t.item_id
	`, userId, start, end).Scan(&list).Error

	return list, err
}


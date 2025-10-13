package transaction

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(transaction *TransactionInsertReq) error {
	tx := &Transaction{
		TransactionID: uuid.New().String(),
		ProductID:    uuid.New().String(),
		Type:        transaction.Type,
		Amount:      transaction.Amount,
		Title:     transaction.Title,
		Price:     transaction.Price,
		Category:  "general",
		Date:      time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(tx)
}

func (s *Service) GetAll() ([]Transaction, error) {
	return s.repo.FindAll()
}

func (s *Service) GetByID(tx_id string, prod_id string) (*Transaction, error) {
	return s.repo.FindByID(tx_id, prod_id)
}

func (s *Service) Update(transaction *Transaction) error {
	return s.repo.Update(transaction)
}

func (s *Service) Delete(tx_id string, prod_id string) error {
	return s.repo.Delete(tx_id, prod_id)
}

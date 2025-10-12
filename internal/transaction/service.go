package transaction

import (
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(transaction *TransactionReq) error {
	tx := &Transaction{
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

func (s *Service) GetByID(id uint) (*Transaction, error) {
	return s.repo.FindByID(id)
}

func (s *Service) Update(transaction *Transaction) error {
	return s.repo.Update(transaction)
}

func (s *Service) Delete(id uint) error {
	return s.repo.Delete(id)
}

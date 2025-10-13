package transaction

import (
	"time"

	r "modjot/internal/response"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) Create(transaction *Transaction) error {
	tx := &Transaction{
		TransactionID: uuid.New().String(),
		ProductID:     uuid.New().String(),
		UserID:        transaction.UserID,
		Type:          transaction.Type,
		Amount:        transaction.Amount,
		Title:         transaction.Title,
		Price:         transaction.Price,
		Category:      transaction.Category,
		Date:          time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repo.Create(tx)
}

func (s *Service) GetAllByUserID(userID string) ([]Transaction, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *Service) GetByID(params *SearchParams) (*Transaction, error) {
	return s.repo.FindByID(params)
}

func (s *Service) Update(transaction *Transaction) error {
	params := &SearchParams{
		TransactionID: transaction.TransactionID,
		ProductID:     transaction.ProductID,
		UserID:        transaction.UserID,
	}
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	if err := validateTransactionOwnership(exists, transaction.UserID); err != nil {
		return err
	}
	return s.repo.Update(transaction)
}

func (s *Service) Delete(params *SearchParams) error {
	return s.repo.Delete(params)
}

func validateTransactionOwnership(tx *Transaction, userID string) error {
	if tx.UserID != userID {
		return r.Conflict(nil, "You are not authorized to access this transaction")
	}
	return nil
}
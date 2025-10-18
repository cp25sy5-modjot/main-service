package transaction

import (
	"time"

	r "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
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

func (s *Service) Update(params *SearchParams, transaction *TransactionUpdateReq) error {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	if err := validateTransactionOwnership(exists, params.UserID); err != nil {
		return err
	}
	_ = utils.MapNonNilStructs(transaction, exists)
	exists.UpdatedAt = time.Now()

	return s.repo.Update(exists)
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

package transactionsvc

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
)

type Service interface {
	GetAllByUserID(userID string) ([]e.TransactionItem, error)
	GetByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error)
	Update(params *m.TransactionItemSearchParams, input *TransactionItemUpdateInput) (*e.TransactionItem, error)
	Delete(params *m.TransactionItemSearchParams) error
}

// concrete implementation
type service struct {
	repo *txirepo.Repository
}

func NewService(repo *txirepo.Repository) *service {
	return &service{repo: repo}
}

func (s *service) GetAllByUserID(userID string) ([]e.TransactionItem, error) {
	transactions, err := s.repo.FindAllByUserID(&m.TransactionItemSearchParams{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *service) GetByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error) {
	tx, err := s.repo.FindByIDWithCategory(params)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *service) Update(params *m.TransactionItemSearchParams, input *TransactionItemUpdateInput) (*e.TransactionItem, error) {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}

	err = utils.MapStructs(input, exists)
	if err != nil {
		return nil, err
	}

	updatedTx, err := s.repo.Update(exists)
	if err != nil {
		return nil, err
	}
	return updatedTx, nil
}

func (s *service) Delete(params *m.TransactionItemSearchParams) error {
	_, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	return s.repo.Delete(params)
}

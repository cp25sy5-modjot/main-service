package transactionsvc

import (
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
	"gorm.io/gorm"
)

type Service interface {
	GetAllByUserID(userID string) ([]e.TransactionItem, error)
	GetByID(params *m.TransactionItemSearchParams) (*e.TransactionItem, error)
	Update(params *m.TransactionItemSearchParams, input *TransactionItemUpdateInput) (*e.TransactionItem, error)
	Delete(params *m.TransactionItemSearchParams) error
}

// concrete implementation
type service struct {
	repo   *txirepo.Repository
	txRepo *txrepo.Repository
}

func NewService(repo *txirepo.Repository, txRepo *txrepo.Repository) *service {
	return &service{repo: repo, txRepo: txRepo}
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
	return s.repo.WithTx(func(tx *gorm.DB) error {

		itemRepo := txirepo.NewRepository(tx)
		transactionRepo := txrepo.NewRepository(tx)

		// check exists
		_, err := itemRepo.FindByID(params)
		if err != nil {
			return err
		}

		// delete item
		if err := itemRepo.DeleteItem(params); err != nil {
			return err
		}

		// count remaining items
		count, err := itemRepo.CountItemsByTransactionID(params.TransactionID)
		if err != nil {
			return err
		}

		// delete transaction if empty
		if count == 0 {
			if err := transactionRepo.Delete(
				&m.TransactionSearchParams{
					TransactionID: params.TransactionID,
					UserID:        params.UserID,
				}); err != nil {
				return err
			}
		}

		return nil
	})
}

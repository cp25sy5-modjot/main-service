package transaction

import (
	"time"

	catRepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	tranRepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"
	"github.com/google/uuid"
)

type Service struct {
	repo     *tranRepo.Repository
	catRepo  *catRepo.Repository
	aiClient pb.AiWrapperServiceClient
}

func NewService(repo *tranRepo.Repository, catRepo *catRepo.Repository, aiClient pb.AiWrapperServiceClient) *Service {
	return &Service{repo: repo, catRepo: catRepo, aiClient: aiClient}
}

func (s *Service) Create(transaction *e.Transaction) (*m.TransactionRes, error) {
	txId := uuid.New().String()
	transaction.Type = "manual"
	tx := buildTransactionObjectToCreate(txId, transaction)
	cat, err := checkCategory(s, tx)
	if err != nil {
		return nil, err
	}
	newTx, err := s.repo.Create(tx)
	if err != nil {
		return nil, err
	}

	return buildTransactionResponse(newTx, cat), nil
}

func (s *Service) ProcessUploadedFile(fileData []byte, userID string) (*m.TransactionRes, error) {
	//fetch user categories
	categories, err := s.catRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	tResponse, err := callAIServiceToBuildTransaction(fileData, categories, s.aiClient)
	if err != nil {
		return nil, err
	}

	return processTransaction(tResponse, categories, userID, s)
}

func (s *Service) GetAllByUserID(userID string) ([]m.TransactionRes, error) {
	transactions, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var transactionResponses []m.TransactionRes
	for _, tx := range transactions {
		cat, err := checkCategory(s, &tx)
		if err != nil {
			return nil, err
		}
		txRes := buildTransactionResponse(&tx, cat)
		transactionResponses = append(transactionResponses, *txRes)
	}
	return transactionResponses, nil
}

func (s *Service) GetAllByUserIDWithFilter(userID string, filter *m.TransactionFilter) ([]m.TransactionRes, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	transactions, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)
	if err != nil {
		return nil, err
	}
	if transactions == nil {
		return []m.TransactionRes{}, nil
	}
	var transactionResponses []m.TransactionRes
	for _, tx := range transactions {
		cat, err := checkCategory(s, &tx)
		if err != nil {
			return nil, err
		}
		txRes := buildTransactionResponse(&tx, cat)
		transactionResponses = append(transactionResponses, *txRes)
	}
	return transactionResponses, nil
}

func (s *Service) GetByID(params *m.TransactionSearchParams) (*m.TransactionRes, error) {
	tx, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	cat, err := checkCategory(s, tx)
	if err != nil {
		return nil, err
	}
	return buildTransactionResponse(tx, cat), nil
}

func (s *Service) Update(params *m.TransactionSearchParams, transaction *m.TransactionUpdateReq) (*m.TransactionRes, error) {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	if err := validateTransactionOwnership(exists, params.UserID); err != nil {
		return nil, err
	}
	_ = utils.MapStructs(transaction, exists)
	cat, err := checkCategory(s, exists)
	if err != nil {
		return nil, err
	}
	updatedTx, err := s.repo.Update(exists)
	if err != nil {
		return nil, err
	}
	return buildTransactionResponse(updatedTx, cat), nil
}

func (s *Service) Delete(params *m.TransactionSearchParams) error {
	return s.repo.Delete(params)
}

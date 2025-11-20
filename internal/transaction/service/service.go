package transaction

import (
	"log"
	"time"

	"context"

	catModel "github.com/cp25sy5-modjot/main-service/internal/category/model"
	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
	r "github.com/cp25sy5-modjot/main-service/internal/response/error"
	tranModel "github.com/cp25sy5-modjot/main-service/internal/transaction/model"
	tranRepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"
	"github.com/google/uuid"
)

type Service struct {
	repo     *tranRepo.Repository
	cat      *catSvc.Service
	aiClient pb.AiWrapperServiceClient
}

func NewService(repo *tranRepo.Repository, cat *catSvc.Service, aiClient pb.AiWrapperServiceClient) *Service {
	return &Service{repo: repo, cat: cat, aiClient: aiClient}
}

func (s *Service) Create(transaction *tranModel.Transaction) (*tranModel.Transaction, error) {
	txId := uuid.New().String()
	transaction.Type = "manual"
	tx := buildTransactionObjectToCreate(txId, transaction)
	if err := checkCategory(s, tx); err != nil {
		return nil, err
	}
	return s.repo.Create(tx)
}

func (s *Service) ProcessUploadedFile(fileData []byte, userID string) (*tranModel.Transaction, error) {
	categoryNames, err := GetCategoryNames(s, userID)
	if err != nil {
		return nil, err
	}

	req := &pb.BuildTransactionFromImageRequest{
		ImageData:  fileData,
		Categories: categoryNames,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // 30 sec timeout for upload
	defer cancel()

	tResponse, err := s.aiClient.BuildTransactionFromImage(ctx, req)
	if err != nil {
		return nil, err
	}
	transaction := &tranModel.Transaction{}
	utils.MapStructs(tResponse, transaction)
	transaction.UserID = userID
	txId := uuid.New().String()
	transaction.Type = "image_upload"
	tx := buildTransactionObjectToCreate(txId, transaction)

	if err := checkCategory(s, tx); err != nil {
		return nil, err
	}

	return s.repo.Create(tx)
}

func (s *Service) GetAllByUserID(userID string) ([]tranModel.Transaction, error) {
	return s.repo.FindAllByUserID(userID)
}

func (s *Service) GetAllByUserIDWithFilter(userID string, filter *tranModel.TransactionFilter) ([]tranModel.Transaction, error) {
	log.Printf("date is %v", filter.Date)
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	transactions, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) GetByID(params *tranModel.TransactionSearchParams) (*tranModel.Transaction, error) {
	return s.repo.FindByID(params)
}

func (s *Service) Update(params *tranModel.TransactionSearchParams, transaction *tranModel.TransactionUpdateReq) error {
	exists, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	if err := validateTransactionOwnership(exists, params.UserID); err != nil {
		return err
	}
	_ = utils.MapStructs(transaction, exists)
	if err := checkCategory(s, exists); err != nil {
		return err
	}
	return s.repo.Update(exists)
}

func (s *Service) Delete(params *tranModel.TransactionSearchParams) error {
	return s.repo.Delete(params)
}

func validateTransactionOwnership(tx *tranModel.Transaction, userID string) error {
	if tx.UserID != userID {
		return r.Conflict(nil, "You are not authorized to access this transaction")
	}
	return nil
}

func GetCategoryNames(s *Service, userID string) ([]string, error) {
	categories, err := s.cat.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	//parse categories to string slice
	var categoryNames []string
	for _, cate := range categories {
		categoryNames = append(categoryNames, cate.CategoryName)
	}
	return categoryNames, nil
}

func buildTransactionObjectToCreate(txId string, tx *tranModel.Transaction) *tranModel.Transaction {
	return &tranModel.Transaction{
		TransactionID: txId,
		ItemID:        uuid.New().String(),
		UserID:        tx.UserID,
		Type:          tx.Type,
		Quantity:      tx.Quantity,
		Title:         tx.Title,
		Price:         tx.Price,
		CategoryID:    tx.CategoryID,
		Date:          time.Now(),
	}
}

func checkCategory(s *Service, tx *tranModel.Transaction) error {
	catSearchParam := &catModel.CategorySearchParams{
		CategoryID: tx.CategoryID,
		UserID:     tx.UserID,
	}
	_, err := s.cat.GetByID(catSearchParam)
	if err != nil {
		return err
	}
	return nil
}

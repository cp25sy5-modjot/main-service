package transaction

import (
	"errors"
	"log"
	"time"

	"context"

	catModel "github.com/cp25sy5-modjot/main-service/internal/category/model"
	catSvc "github.com/cp25sy5-modjot/main-service/internal/category/service"
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

func (s *Service) Create(transaction *tranModel.Transaction) (*tranModel.TransactionRes, error) {
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

func (s *Service) ProcessUploadedFile(fileData []byte, userID string) (*tranModel.TransactionRes, error) {
	//fetch user categories
	categories, err := s.cat.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	tResponse, err := callAIServiceToBuildTransaction(fileData, categories, s.aiClient)
	if err != nil {
		return nil, err
	}

	return processTransaction(tResponse, categories, userID, s)
}

func (s *Service) GetAllByUserID(userID string) ([]tranModel.TransactionRes, error) {
	transactions, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var transactionResponses []tranModel.TransactionRes
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

func (s *Service) GetAllByUserIDWithFilter(userID string, filter *tranModel.TransactionFilter) ([]tranModel.TransactionRes, error) {
	log.Printf("date is %v", filter.Date)
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	transactions, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)
	if err != nil {
		return nil, err
	}
	if transactions == nil {
		return []tranModel.TransactionRes{}, nil
	}
	var transactionResponses []tranModel.TransactionRes
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

func (s *Service) GetByID(params *tranModel.TransactionSearchParams) (*tranModel.TransactionRes, error) {
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

func (s *Service) Update(params *tranModel.TransactionSearchParams, transaction *tranModel.TransactionUpdateReq) (*tranModel.TransactionRes, error) {
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

func (s *Service) Delete(params *tranModel.TransactionSearchParams) error {
	return s.repo.Delete(params)
}

func validateTransactionOwnership(tx *tranModel.Transaction, userID string) error {
	if tx.UserID != userID {
		return errors.New("you are not authorized to access this transaction")
	}
	return nil
}

func GetCategoryNames(categories []catModel.Category) ([]string, error) {
	//parse categories to string slice
	var categoryNames []string
	for _, cate := range categories {
		categoryNames = append(categoryNames, cate.CategoryName)
	}
	return categoryNames, nil
}

func checkCategory(s *Service, tx *tranModel.Transaction) (*catModel.Category, error) {
	catSearchParam := &catModel.CategorySearchParams{
		CategoryID: tx.CategoryID,
		UserID:     tx.UserID,
	}
	cat, err := s.cat.GetByID(catSearchParam)
	if err != nil {
		return nil, errors.New("category does not exist")
	}
	return cat, nil
}

func buildTransactionObjectToCreate(txId string, tx *tranModel.Transaction) *tranModel.Transaction {
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	return &tranModel.Transaction{
		TransactionID: txId,
		ItemID:        uuid.New().String(),
		UserID:        tx.UserID,
		Type:          tx.Type,
		Quantity:      tx.Quantity,
		Title:         tx.Title,
		Price:         tx.Price,
		CategoryID:    tx.CategoryID,
		Date:          tx.Date,
	}
}

func buildTransactionResponse(tx *tranModel.Transaction, category *catModel.Category) *tranModel.TransactionRes {
	return &tranModel.TransactionRes{
		TransactionID:     tx.TransactionID,
		ItemID:            tx.ItemID,
		Type:              tx.Type,
		Title:             tx.Title,
		Price:             tx.Price,
		Quantity:          tx.Quantity,
		TotalPrice:        tx.Price * float64(tx.Quantity),
		Date:              tx.Date,
		CategoryID:        tx.CategoryID,
		CategoryName:      category.CategoryName,
		CategoryColorCode: category.ColorCode,
	}
}

func matchCategoryFromName(categories []catModel.Category, categoryName string) *catModel.Category {
	for _, cat := range categories {
		if cat.CategoryName == categoryName {
			return &cat
		}
	}
	return nil
}

func callAIServiceToBuildTransaction(fileData []byte, categories []catModel.Category, aiClient pb.AiWrapperServiceClient) (*pb.TransactionResponse, error) {
	//get category names to send to ai service
	categoryNames, err := GetCategoryNames(categories)
	if err != nil {
		return nil, err
	}
	req := &pb.BuildTransactionFromImageRequest{
		ImageData:  fileData,
		Categories: categoryNames,
	}
	const timeout = 5*time.Minute + 30*time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout) // 30 sec timeout for upload
	defer cancel()

	tResponse, err := aiClient.BuildTransactionFromImage(ctx, req)
	if err != nil {
		return nil, err
	}
	return tResponse, nil
}

func processTransaction(tResponse *pb.TransactionResponse, categories []catModel.Category, userID string, s *Service) (*tranModel.TransactionRes, error) {
	match := matchCategoryFromName(categories, tResponse.Category)
	if match == nil {
		return nil, errors.New("category does not exist")
	}
	transaction := &tranModel.Transaction{}
	utils.MapStructs(tResponse, transaction)
	transaction.UserID = userID
	transaction.CategoryID = match.CategoryID
	txId := uuid.New().String()
	transaction.Type = "image_upload"
	tx := buildTransactionObjectToCreate(txId, transaction)
	newTx, err := s.repo.Create(tx)
	if err != nil {
		return nil, err
	}
	log.Printf("Created transaction: %+v", newTx)
	return buildTransactionResponse(newTx, match), nil
}

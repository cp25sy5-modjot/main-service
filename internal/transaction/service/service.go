package transactionsvc

import (
	"context"
	"errors"
	"log"
	"time"

	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"
	"github.com/google/uuid"
)

type Service struct {
	repo     *txrepo.Repository
	catrepo  *catrepo.Repository
	aiClient pb.AiWrapperServiceClient
}

func NewService(repo *txrepo.Repository, catrepo *catrepo.Repository, aiClient pb.AiWrapperServiceClient) *Service {
	return &Service{repo: repo, catrepo: catrepo, aiClient: aiClient}
}

func (s *Service) Create(userID string, input *TransactionCreateInput) (*e.Transaction, error) {
	txId := uuid.New().String()

	_, err := s.catrepo.FindByID(&m.CategorySearchParams{
		CategoryID: input.CategoryID,
		UserID:     userID,
	})
	if err != nil {
		return nil, err
	}

	tx := buildTransactionObjectToCreate(txId, userID, "manual", input)
	txWithCat, err := saveNewTransaction(s, tx)
	if err != nil {
		return nil, err
	}
	return txWithCat, nil
}

func (s *Service) ProcessUploadedFile(fileData []byte, userID string) (*e.Transaction, error) {
	if s.aiClient == nil {
		return nil, errors.New("AI client not configured (this method should only be used in worker process)")
	}

	// 1. fetch categories
	categories, err := s.catrepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 2. call AI service
	resp, err := callAIServiceToBuildTransaction(fileData, categories, s.aiClient)
	if err != nil {
		return nil, err
	}

	log.Printf("AI Service Response: %+v", resp)

	// 3. process into real transaction (same as before)
	return processTransaction(resp, categories, userID, s)
}

func (s *Service) GetAllByUserID(userID string) ([]e.Transaction, error) {
	transactions, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) GetAllByUserIDWithFilter(userID string, filter *m.TransactionFilter) ([]e.Transaction, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	filter.PreviousMonth = false
	transactions, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)

	if err != nil {
		return nil, err
	}

	if transactions == nil {
		return []e.Transaction{}, nil
	}

	return transactions, nil
}

type MonthlyResult struct {
	CurrentMonth  []e.Transaction `json:"current_month"`
	PreviousMonth []e.Transaction `json:"previous_month"`
}

func (s *Service) GetAllComparePreviousMonthAndByUserIDWithFilter(userID string, filter *m.TransactionFilter) (*MonthlyResult, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}

	// --- Current Month ---
	filter.PreviousMonth = false
	current, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)
	if err != nil {
		return nil, err
	}

	// --- Previous Month ---
	filter.PreviousMonth = true
	previous, err := s.repo.FindAllByUserIDAndFiltered(userID, filter)
	if err != nil {
		return nil, err
	}

	return &MonthlyResult{
		CurrentMonth:  current,
		PreviousMonth: previous,
	}, nil
}

func (s *Service) GetByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	tx, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Service) Update(params *m.TransactionSearchParams, input *TransactionUpdateInput) (*e.Transaction, error) {
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

func (s *Service) Delete(params *m.TransactionSearchParams) error {
	_, err := s.repo.FindByID(params)
	if err != nil {
		return err
	}
	return s.repo.Delete(params)
}

// utils functions for service
func GetCategoryNames(categories []e.Category) ([]string, error) {
	//parse categories to string slice
	var categoryNames []string
	for _, cate := range categories {
		categoryNames = append(categoryNames, cate.CategoryName)
	}
	return categoryNames, nil
}

func buildTransactionObjectToCreate(txId, userID, txType string, tx *TransactionCreateInput) *e.Transaction {
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	return &e.Transaction{
		TransactionID: txId,
		ItemID:        uuid.New().String(),
		UserID:        userID,
		Type:          txType,
		Quantity:      tx.Quantity,
		Title:         tx.Title,
		Price:         tx.Price,
		CategoryID:    tx.CategoryID,
		Date:          tx.Date,
	}
}

// func buildTransactionObjectToCreates(txId, userID, txType string,  txs []*TransactionCreateInput) []*e.Transaction {
// 	var transactions []*e.Transaction
// 	for _, tx := range txs {
// 		newTx := buildTransactionObjectToCreate(txId, userID, txType, tx)
// 		transactions = append(transactions, newTx)
// 	}
// 	return transactions
// }

func matchCategoryFromName(categories []e.Category, categoryName string) *e.Category {
	for _, cat := range categories {
		if cat.CategoryName == categoryName {
			return &cat
		}
	}
	return nil
}

func callAIServiceToBuildTransaction(fileData []byte, categories []e.Category, aiClient pb.AiWrapperServiceClient) (*pb.TransactionResponse, error) {
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

func processTransaction(tResponse *pb.TransactionResponse, categories []e.Category, userID string, s *Service) (*e.Transaction, error) {
	match := matchCategoryFromName(categories, tResponse.Category)
	if match == nil {
		return nil, errors.New("category does not exist")
	}
	transaction := &TransactionCreateInput{}
	err := utils.MapStructs(tResponse, transaction)
	if err != nil {
		return nil, err
	}
	transaction.CategoryID = &match.CategoryID
	txId := uuid.New().String()

	tx := buildTransactionObjectToCreate(txId, userID, "image_upload", transaction)
	txWithCat, err := saveNewTransaction(s, tx)
	if err != nil {
		return nil, err
	}
	return txWithCat, nil
}

func saveNewTransaction(s *Service, tx *e.Transaction) (*e.Transaction, error) {
	newTx, err := s.repo.Create(tx)
	if err != nil {
		return nil, err
	}
	// Reload with preload
	txWithCat, err := s.repo.FindByID(&m.TransactionSearchParams{
		TransactionID: newTx.TransactionID,
		ItemID:        newTx.ItemID,
		UserID:        newTx.UserID,
	})
	if err != nil {
		return nil, err
	}
	return txWithCat, nil
}

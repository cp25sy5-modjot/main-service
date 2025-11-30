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

type Service interface {
	Create(userID string, input *TransactionCreateInput) (*e.Transaction, error)
	ProcessUploadedFile(fileData []byte, userID string) (*e.Transaction, error)

	GetAllByUserID(userID string) ([]e.Transaction, error)
	GetAllByUserIDWithFilter(userID string, filter *m.TransactionFilter) ([]e.Transaction, error)
	GetAllComparePreviousMonthAndByUserIDWithFilter(userID string, filter *m.TransactionFilter) (*MonthlyResult, error)

	GetByID(params *m.TransactionSearchParams) (*e.Transaction, error)
	Update(params *m.TransactionSearchParams, input *TransactionUpdateInput) (*e.Transaction, error)
	Delete(params *m.TransactionSearchParams) error
}

// concrete implementation
type service struct {
	repo     *txrepo.Repository
	catrepo  *catrepo.Repository
	aiClient pb.AiWrapperServiceClient
}

func NewService(repo *txrepo.Repository, catrepo *catrepo.Repository, aiClient pb.AiWrapperServiceClient) *service {
	return &service{repo: repo, catrepo: catrepo, aiClient: aiClient}
}

func (s *service) Create(userID string, input *TransactionCreateInput) (*e.Transaction, error) {
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

func (s *service) ProcessUploadedFile(fileData []byte, userID string) (*e.Transaction, error) {
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

func (s *service) GetAllByUserID(userID string) ([]e.Transaction, error) {
	transactions, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *service) GetAllByUserIDWithFilter(userID string, filter *m.TransactionFilter) ([]e.Transaction, error) {
	if filter.Date == nil {
		now := time.Now()
		filter.Date = &now
	}
	start, end := utils.GetStartAndEndOfMonth(*filter.Date)
	transactions, err := s.repo.FindAllByUserIDAndFiltered(userID, start, end)

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

func (s *service) GetAllComparePreviousMonthAndByUserIDWithFilter(userID string, filter *m.TransactionFilter) (*MonthlyResult, error) {
	// --- Current Month ---
	start, end := utils.GetStartAndEndOfMonth(*filter.Date)
	current, err := s.repo.FindAllByUserIDAndFiltered(userID, start, end)
	if err != nil {
		return nil, err
	}

	// --- Previous Month ---
	date := filter.Date.AddDate(0, -1, 0) // move filter date to previous month
	previousStart, previousEnd := utils.GetStartAndEndOfMonth(date)
	previous, err := s.repo.FindAllByUserIDAndFiltered(userID, previousStart, previousEnd)
	if err != nil {
		return nil, err
	}

	return &MonthlyResult{
		CurrentMonth:  current,
		PreviousMonth: previous,
	}, nil
}

func (s *service) GetByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	tx, err := s.repo.FindByID(params)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *service) Update(params *m.TransactionSearchParams, input *TransactionUpdateInput) (*e.Transaction, error) {
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

func (s *service) Delete(params *m.TransactionSearchParams) error {
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

func isDefaultDate(t time.Time) bool {
	return t.Year() == 1 && t.Month() == time.January && t.Day() == 1
}

func buildTransactionObjectToCreate(txId, userID, txType string, tx *TransactionCreateInput) *e.Transaction {
	if isDefaultDate(tx.Date) {
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

func processTransaction(tResponse *pb.TransactionResponse, categories []e.Category, userID string, s *service) (*e.Transaction, error) {
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

func saveNewTransaction(s *service, tx *e.Transaction) (*e.Transaction, error) {
	if tx.Quantity != 1 {
		//quantity will be remove in future, so we multiply price with quantity and set quantity to 1
		tx.Price *= float64(tx.Quantity)
		tx.Quantity = 1
	}
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

func getMonthRange(filter *m.TransactionFilter) (time.Time, time.Time) {

	t := filter.Date
	loc := t.Location()
	// First day of current month at 00:00
	firstOfCurrent := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)

	// Last day of previous month = firstOfCurrent - 1 day
	lastOfPrevious := firstOfCurrent.AddDate(0, 0, -1)

	// Last day of current month = first of next month - 1 day
	firstOfNext := firstOfCurrent.AddDate(0, 1, 0)
	lastOfCurrent := firstOfNext.AddDate(0, 0, -1)

	// Set final times to 17:00:00
	startRange := time.Date(
		lastOfPrevious.Year(), lastOfPrevious.Month(), lastOfPrevious.Day(),
		17, 0, 0, 0,
		loc,
	)

	endRange := time.Date(
		lastOfCurrent.Year(), lastOfCurrent.Month(), lastOfCurrent.Day(),
		17, 0, 0, 0,
		loc,
	)
	return startRange, endRange
}

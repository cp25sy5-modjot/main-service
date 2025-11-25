package transaction

import (
	"context"
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"
	"github.com/google/uuid"
)

func validateTransactionOwnership(tx *e.Transaction, userID string) error {
	if tx.UserID != userID {
		return errors.New("you are not authorized to access this transaction")
	}
	return nil
}

func GetCategoryNames(categories []e.Category) ([]string, error) {
	//parse categories to string slice
	var categoryNames []string
	for _, cate := range categories {
		categoryNames = append(categoryNames, cate.CategoryName)
	}
	return categoryNames, nil
}

func checkCategory(s *Service, tx *e.Transaction) (*e.Category, error) {
	catSearchParam := &m.CategorySearchParams{
		CategoryID: tx.CategoryID,
		UserID:     tx.UserID,
	}
	cat, err := s.catRepo.FindByID(catSearchParam)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}
	return cat, nil
}

func buildTransactionObjectToCreate(txId string, tx *e.Transaction) *e.Transaction {
	if tx.Date.IsZero() {
		tx.Date = time.Now()
	}
	return &e.Transaction{
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

func buildTransactionResponse(tx *e.Transaction, category *e.Category) *m.TransactionRes {
	return &m.TransactionRes{
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

func processTransaction(tResponse *pb.TransactionResponse, categories []e.Category, userID string, s *Service) (*m.TransactionRes, error) {
	match := matchCategoryFromName(categories, tResponse.Category)
	if match == nil {
		return nil, errors.New("category does not exist")
	}
	transaction := &e.Transaction{}
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
	return buildTransactionResponse(newTx, match), nil
}

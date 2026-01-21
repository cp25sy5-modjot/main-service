package transactionsvc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	"github.com/cp25sy5-modjot/main-service/internal/shared/utils"
	draft "github.com/cp25sy5-modjot/main-service/internal/draft"
	txrepo "github.com/cp25sy5-modjot/main-service/internal/transaction/repository"
	txirepo "github.com/cp25sy5-modjot/main-service/internal/transaction_item/repository"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Create(userID string, input *m.TransactionCreateInput) (*e.Transaction, error)
	ProcessUploadedFile(fileData []byte, userID string) (*draft.DraftTxn, error)

	GetAllByUserID(userID string) ([]e.Transaction, error)
	GetAllByUserIDWithFilter(userID string, filter *m.TransactionFilter) ([]e.Transaction, error)
	GetAllComparePreviousMonthAndByUserIDWithFilter(userID string, filter *m.TransactionFilter) (*MonthlyResult, error)

	GetByID(params *m.TransactionSearchParams) (*e.Transaction, error)
	Update(params *m.TransactionSearchParams, input *m.TransactionUpdateInput) (*e.Transaction, error)
	Delete(params *m.TransactionSearchParams) error
}

// concrete implementation
type service struct {
	db       *gorm.DB
	repo     *txrepo.Repository
	txirepo  *txirepo.Repository
	catrepo  *catrepo.Repository

	draftRepo *draft.DraftRepository

	aiClient pb.AiWrapperServiceClient
}

func NewService(db *gorm.DB, repo *txrepo.Repository, txirepo *txirepo.Repository, catrepo *catrepo.Repository, aiClient pb.AiWrapperServiceClient) *service {
	return &service{db: db, repo: repo, txirepo: txirepo, catrepo: catrepo, aiClient: aiClient}
}

func (s *service) Create(
	userID string,
	input *m.TransactionCreateInput,
) (*e.Transaction, error) {
	return s.CreateInternal(userID, e.TransactionManual, input)
}

func (s *service) CreateInternal(
	userID string,
	txType e.TransactionType,
	input *m.TransactionCreateInput,
) (*e.Transaction, error) {

	// 1. validate
	if err := s.validateCreateInput(userID, input); err != nil {
		return nil, err
	}

	// 2. build
	txID := uuid.New().String()
	tx, items := buildTransactionToCreate(txID, userID, txType, input)

	// 3. save (atomic)
	return s.saveNewTransaction(tx, items)
}

func (s *service) ProcessUploadedFile(
	fileData []byte,
	userID string,
) (*draft.DraftTxn, error) {

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

	return mapToDraft(resp, categories, userID)
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
		now := time.Now().UTC()
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
	current, err := s.repo.FindAllByUserIDWithRelationsAndFiltered(userID, start, end)
	if err != nil {
		return nil, err
	}

	// --- Previous Month ---
	previousStart, previousEnd := utils.GetStartAndEndOfPreviousMonth(*filter.Date)
	previous, err := s.repo.FindAllByUserIDWithRelationsAndFiltered(userID, previousStart, previousEnd)
	if err != nil {
		return nil, err
	}

	return &MonthlyResult{
		CurrentMonth:  current,
		PreviousMonth: previous,
	}, nil
}

func (s *service) GetByID(params *m.TransactionSearchParams) (*e.Transaction, error) {
	tx, err := s.repo.FindByIDWithRelations(params)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *service) Update(
	params *m.TransactionSearchParams,
	input *m.TransactionUpdateInput,
) (*e.Transaction, error) {

	exists, err := s.repo.FindByIDWithRelations(params)
	if err != nil {
		return nil, err
	}

	// --- Date ---
	if input.Date != nil {
		utc := input.Date.UTC()
		exists.Date = utc
	}

	// --- Items (PATCH replace semantics) ---
	if input.Items != nil {
		err := s.txirepo.DeleteByTransactionID(exists.TransactionID)
		if err != nil {
			return nil, err
		}

		newItems, err := ReplaceTransactionItems(
			exists.TransactionID,
			input.Items,
		)
		if err != nil {
			return nil, err
		}

		if err := s.txirepo.CreateMany(newItems); err != nil {
			return nil, err
		}
	}

	return s.repo.FindByIDWithRelations(params)
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

func buildTransactionToCreate(
	txID, userID string,
	txType e.TransactionType,
	input *m.TransactionCreateInput,
) (*e.Transaction, []e.TransactionItem) {

	if isDefaultDate(input.Date) {
		input.Date = time.Now().UTC()
	}

	items := make([]e.TransactionItem, 0, len(input.Items))
	for _, it := range input.Items {
		items = append(items, e.TransactionItem{
			TransactionID: txID,
			ItemID:        uuid.New().String(),
			Title:         it.Title,
			Price:         it.Price,
			CategoryID:    it.CategoryID,
		})
	}

	tx := &e.Transaction{
		TransactionID: txID,
		UserID:        userID,
		Type:          txType,
		Date:          input.Date.UTC(),
	}

	return tx, items
}

func matchCategoryFromName(categories []e.Category, categoryName string) *e.Category {
	for i := range categories {
		if categories[i].CategoryName == categoryName {
			return &categories[i]
		}
	}
	return nil
}

func callAIServiceToBuildTransaction(fileData []byte, categories []e.Category, aiClient pb.AiWrapperServiceClient) (*pb.TransactionResponseV2, error) {
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

func (s *service) saveNewTransaction(
	tx *e.Transaction,
	items []e.TransactionItem,
) (*e.Transaction, error) {

	if tx == nil {
		return nil, errors.New("transaction is nil")
	}
	if len(items) == 0 {
		return nil, errors.New("no transaction items to save")
	}

	err := s.db.Transaction(func(db *gorm.DB) error {
		if err := s.repo.WithTx(db).Create(tx); err != nil {
			return err
		}

		if err := s.txirepo.WithTx(db).CreateMany(items); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return s.repo.FindByIDWithRelations(&m.TransactionSearchParams{
		TransactionID: tx.TransactionID,
		UserID:        tx.UserID,
	})
}

func (s *service) validateCreateInput(
	userID string,
	input *m.TransactionCreateInput,
) error {

	if input == nil {
		return errors.New("input is required")
	}

	if len(input.Items) == 0 {
		return errors.New("at least one item is required")
	}
	categories, err := s.catrepo.FindAllByUserID(userID)
	if err != nil {
		return err
	}

	categoryMap := map[string]bool{}
	for _, c := range categories {
		categoryMap[c.CategoryID] = true
	}

	for _, it := range input.Items {
		if it.Title == "" {
			return errors.New("item title is required")
		}
		if it.Price < 0 {
			return errors.New("item price must be positive")
		}
		if !categoryMap[it.CategoryID] {
			return errors.New("invalid category")
		}
	}

	return nil
}

func ReplaceTransactionItems(
	txID string,
	input []m.TransactionItemInput,
) ([]e.TransactionItem, error) {

	if len(input) == 0 {
		return nil, errors.New("items cannot be empty")
	}

	items := make([]e.TransactionItem, 0, len(input))

	for _, it := range input {
		items = append(items, e.TransactionItem{
			TransactionID: txID,
			ItemID:        uuid.New().String(),
			Title:         it.Title,
			Price:         it.Price,
			CategoryID:    it.CategoryID,
		})
	}

	return items, nil
}

func mapToDraft(
	resp *pb.TransactionResponseV2,
	categories []e.Category,
	userID string,
) (*draft.DraftTxn, error) {

	if len(resp.Items) == 0 {
		return nil, errors.New("no transaction items")
	}

	// parse date
	date := time.Now().UTC()
	if resp.Date != "" {
		if parsed, err := time.Parse(time.RFC3339, resp.Date); err == nil {
			date = parsed.UTC()
		}
	}

	var items []draft.DraftItem

	for _, res := range resp.Items {

		match := matchCategoryFromName(categories, res.Category)
		if match == nil {
			return nil, fmt.Errorf("category not found: %s", res.Category)
		}

		items = append(items, draft.DraftItem{
			Title:      res.Title,
			Price:      res.Price,
			CategoryID: match.CategoryID,
		})
	}

	return &draft.DraftTxn{
		UserID: userID,
		Status: draft.DraftStatusWaitingConfirm,

		Date:  date,
		Items: items,

		CreatedAt: time.Now(),
	}, nil
}

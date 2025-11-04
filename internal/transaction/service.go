package transaction

import (
	"time"

	"context"

	r "github.com/cp25sy5-modjot/main-service/internal/response/error"
	"github.com/cp25sy5-modjot/main-service/internal/utils"
	pb "github.com/cp25sy5-modjot/proto/gen/ai/v1"
	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	aiClient pb.AiWrapperServiceClient
}

func NewService(repo *Repository, aiClient pb.AiWrapperServiceClient) *Service {
	return &Service{repo, aiClient}
}

func (s *Service) Create(transaction *Transaction) error {
	tx := &Transaction{
		TransactionID: uuid.New().String(),
		ItemID:        uuid.New().String(),
		UserID:        transaction.UserID,
		Type:          transaction.Type,
		Quantity:      transaction.Quantity,
		Title:         transaction.Title,
		Price:         transaction.Price,
		Category:      transaction.Category,
		Date:          time.Now(),
	}
	return s.repo.Create(tx)
}

func (s *Service) ProcessUploadedFile(fileData []byte, userID string) error {

	req := &pb.BuildTransactionFromImageRequest{
		ImageData:  fileData,
		Categories: []string{"food", "transportation", "utilities", "entertainment", "health", "other"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15 sec timeout for upload
	defer cancel()

	tResponse, err := s.aiClient.BuildTransactionFromImage(ctx, req)
	if err != nil {
		return err
	}
	tx := &Transaction{}
	utils.MapNonNilStructs(tResponse, tx)

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

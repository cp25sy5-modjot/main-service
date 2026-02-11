package draft

import (
	"context"
	"errors"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	catrepo "github.com/cp25sy5-modjot/main-service/internal/category/repository"
)

type Service interface {

	GetDraft(ctx context.Context, traceID, userID string) (*DraftTxn, error)
	ListDraft(ctx context.Context, userID string) ([]DraftTxn, error)
	SaveDraft(ctx context.Context, traceID, userID string, req NewDraftRequest) (*DraftTxn, error)
	UpdateDraft(ctx context.Context, traceID, userID string, req ConfirmRequest) (*DraftTxn, error)
	ConfirmDraft(ctx context.Context, traceID string, userID string, req ConfirmRequest) (*e.Transaction, error)
	DeleteDraft(ctx context.Context, traceID string) error
	GetDraftStats(ctx context.Context, userID string) (*DraftStats, error)
	GetDraftWithCategory(ctx context.Context, traceID, userID string) (*DraftRes, error)
	ListDraftWithCategory(ctx context.Context, userID string) ([]DraftRes, error)
}

type service struct {
	draftRepo *DraftRepository
	categoryRepo *catrepo.Repository

	createInternal func(
		userID string,
		transactionType e.TransactionType,
		input *m.TransactionCreateInput,
	) (*e.Transaction, error)
}

func NewService(
	repo *DraftRepository,
	categoryRepo *catrepo.Repository,
	createFn func(string, e.TransactionType, *m.TransactionCreateInput) (*e.Transaction, error),
) Service {
	return &service{
		draftRepo:      repo,
		categoryRepo:   categoryRepo,
		createInternal: createFn,
	}
}

func (s *service) GetDraft(ctx context.Context, traceID, userID string) (*DraftTxn, error) {

	d, err := s.draftRepo.Get(ctx, traceID)
	if err != nil {
		return nil, err
	}

	if d.UserID != userID {
		return nil, errors.New("not owner")
	}

	return d, nil
}

func (s *service) ListDraft(ctx context.Context, userID string) ([]DraftTxn, error) {
	return s.draftRepo.ListByUser(ctx, userID)
}

func (s *service) UpdateDraft(
	ctx context.Context,
	traceID string,
	userID string,
	req ConfirmRequest,
) (*DraftTxn, error) {

	d, err := s.draftRepo.Get(ctx, traceID)
	if err != nil {
		return nil, errors.New("draft not found")
	}

	if d.UserID != userID {
		return nil, errors.New("not owner")
	}

	if d.Status != DraftStatusWaitingConfirm {
		return nil, errors.New("cannot edit at this stage")
	}

	// apply change
	if len(req.Items) > 0 {
		d.Items = req.Items
	}

	if req.Date != nil {
		d.Date = *req.Date
	}

	// validate
	for _, it := range d.Items {
		if it.Price <= 0 {
			return nil, errors.New("price must be > 0")
		}
	}

	d.UpdatedAt = time.Now()

	if err := s.draftRepo.Save(ctx, *d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *service) SaveDraft(
	ctx context.Context,
	traceID string,
	userID string,
	req NewDraftRequest,
) (*DraftTxn, error) {

	for _, it := range req.Items {
		if it.Price <= 0 {
			return nil, errors.New("price must be > 0")
		}
	}

	d := &DraftTxn{
		Title:     req.Title,
		TraceID:   traceID,
		UserID:    userID,
		Status:    DraftStatusProcessing,
		Date:      req.Date,
		Items:     req.Items,
		CreatedAt: req.CreatedAt,
		UpdatedAt: req.CreatedAt,
	}

	if err := s.draftRepo.Save(ctx, *d); err != nil {
		return nil, err
	}

	return d, nil
}

func (s *service) ConfirmDraft(
	ctx context.Context,
	traceID string,
	userID string,
	req ConfirmRequest,
) (*e.Transaction, error) {

	d, err := s.draftRepo.Get(ctx, traceID)
	if err != nil {
		return nil, errors.New("draft not found")
	}

	if d.UserID != userID {
		return nil, errors.New("not owner")
	}

	if d.Status != DraftStatusWaitingConfirm {
		return nil, errors.New("draft not ready")
	}

	if len(req.Items) == 0 {
		return nil, errors.New("cannot confirm empty draft")
	}

	for _, it := range req.Items {
		if it.Price < 0 {
			return nil, errors.New("price must be > 0")
		}
	}

	input := mapConfirmDraftToCreateInput(&req)

	tx, err := s.createInternal(userID, e.TransactionUpload, input)
	if err != nil {
		return nil, err
	}

	_ = s.draftRepo.Delete(ctx, traceID)

	return tx, nil
}

func (s *service) DeleteDraft(ctx context.Context, traceID string) error {
	return s.draftRepo.Delete(ctx, traceID)
}

func (s *service) GetDraftStats(ctx context.Context, userID string) (*DraftStats, error) {
	return s.draftRepo.StatsByUser(ctx, userID)
}

func (s *service) GetDraftWithCategory(
	ctx context.Context,
	traceID, userID string,
) (*DraftRes, error) {

	d, err := s.GetDraft(ctx, traceID, userID)
	if err != nil {
		return nil, err
	}

	ids := uniqueCategoryIDsFromDrafts([]DraftTxn{*d})

	categoryMap, err := s.categoryRepo.FindByIDs(ctx, userID, ids)
	if err != nil {
		return nil, err
	}

	res := buildDraftRes(*d, categoryMap)

	return &res, nil
}

func (s *service) ListDraftWithCategory(
	ctx context.Context,
	userID string,
) ([]DraftRes, error) {

	drafts, err := s.draftRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(drafts) == 0 {
		return []DraftRes{}, nil
	}

	ids := uniqueCategoryIDsFromDrafts(drafts)

	categoryMap, err := s.categoryRepo.FindByIDs(ctx, userID, ids)
	if err != nil {
		return nil, err
	}

	result := make([]DraftRes, 0, len(drafts))

	for _, d := range drafts {
		result = append(result, buildDraftRes(d, categoryMap))
	}

	return result, nil
}

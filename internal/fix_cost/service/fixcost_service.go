package fixcostsvc

import (
	"context"
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type Service interface {
	GetByID(ctx context.Context, id string, userID string) (*e.FixCost, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error)
	Create(ctx context.Context, input *m.FixCostCreateInput) (*e.FixCost, error)
	Update(ctx context.Context, input *m.FixCostUpdateInput) (*e.FixCost, error)
	Delete(ctx context.Context, id string, userID string) error
}

type service struct {
	repo     fcrepo.Repository
	redisOpt asynq.RedisClientOpt
}

func NewService(repo fcrepo.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetByID(ctx context.Context, id string, userID string) (*e.FixCost, error) {
	fc, err := s.repo.FindByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	return fc, nil
}

func (s *service) GetAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error) {
	fcs, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return fcs, nil
}

func (s *service) Create(ctx context.Context, input *m.FixCostCreateInput) (*e.FixCost, error) {

	fcId := uuid.New().String()
	newfc := e.FixCost{
		FixCostID:     fcId,
		UserID:        input.UserID,
		Title:         input.Title,
		Price:         input.Price,
		CategoryID:    input.CategoryID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		NextRunDate:   input.StartDate, // first run is on start date
		MaxRun:        input.MaxRun,
		RunCount:      0,
		IntervalType:  e.IntervalType(input.IntervalType),
		IntervalValue: input.IntervalValue,
		Status:        calculateStatus(input.EndDate, input.MaxRun)}

	err := s.repo.Create(ctx, &newfc)
	if err != nil {
		return nil, err
	}

	fc, err := s.repo.FindByID(ctx, fcId, newfc.UserID)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	if fc.NextRunDate.Equal(today) || fc.NextRunDate.Before(today) {
		
	}

	return fc, nil
}

func (s *service) Update(ctx context.Context, input *m.FixCostUpdateInput) (*e.FixCost, error) {

	exists, err := s.repo.FindByID(ctx, input.FixCostID, input.UserID)
	if err != nil {
		return nil, err
	}

	if exists == nil {
		return nil, gorm.ErrRecordNotFound
	}

	// update fields
	exists.EndDate = input.EndDate
	exists.MaxRun = input.MaxRun

	if input.Title != nil {
		exists.Title = *input.Title
	}

	if input.CategoryID != nil {
		exists.CategoryID = *input.CategoryID
	}

	if input.Price != nil {
		exists.Price = *input.Price
	}

	if input.IntervalType != nil {
		exists.IntervalType = e.IntervalType(*input.IntervalType)
	}

	if input.IntervalValue != nil {
		exists.IntervalValue = *input.IntervalValue
	}

	if input.Status != nil {
		exists.Status = e.FixCostStatus(*input.Status)
	}

	if input.StartDate != nil && !input.StartDate.Equal(exists.StartDate) {
		exists.StartDate = *input.StartDate
		exists.NextRunDate = CalculateNextRun(*exists)
	}

	err = s.repo.Update(ctx, exists)
	if err != nil {
		return nil, err
	}

	fc, err := s.repo.FindByID(ctx, input.FixCostID, input.UserID)
	if err != nil {
		return nil, err
	}

	return fc, nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	// 1. ลบ fixcost
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	inspector := asynq.NewInspector(s.redisOpt)

	_ = inspector.DeleteTask("default", "fixcost:"+id)

	return nil
}

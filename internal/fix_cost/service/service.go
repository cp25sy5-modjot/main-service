package fixcostsvc

import (
	"context"
	"log"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	fcrepo "github.com/cp25sy5-modjot/main-service/internal/fix_cost/repository"
	"github.com/cp25sy5-modjot/main-service/internal/jobs/tasks"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type Service interface {
	GetByID(ctx context.Context, id string, userID string) (*e.FixCost, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error)
	Create(ctx context.Context, input *m.FixCostCreateInput) (*e.FixCost, error)
	Update(ctx context.Context, input *m.FixCostUpdateInput) error
	Delete(ctx context.Context, id string, userID string) error

	RecoverFixCostJobs()
}

type service struct {
	repo *fcrepo.Repository

	asynqClient *asynq.Client
}

func NewService(repo *fcrepo.Repository, asynqClient *asynq.Client) Service {
	return &service{
		repo:        repo,
		asynqClient: asynqClient,
	}
}

func (s *service) GetByID(ctx context.Context, id string, userID string) (*e.FixCost, error) {
	return s.repo.FindByID(ctx, id, userID)
}

func (s *service) GetAllByUserID(ctx context.Context, userID string) ([]*e.FixCost, error) {
	return s.repo.FindAllByUserID(ctx, userID)
}

func (s *service) Create(ctx context.Context, input *m.FixCostCreateInput) (*e.FixCost, error) {

	fc := e.FixCost{
		FixCostID:     uuid.New().String(),
		UserID:        input.UserID,
		Title:         input.Title,
		Price:         input.Price,
		CategoryID:    input.CategoryID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		NextRunDate:   input.StartDate, // first run is on start date
		RemainingRuns: input.RemainingRuns,
		IntervalType:  e.IntervalType(input.IntervalType),
		IntervalValue: input.IntervalValue,
		Status:        e.FixCostStatusActive,
	}

	err := s.repo.Create(ctx, &fc)
	if err != nil {
		return nil, err
	}
	task, _ := tasks.NewRunFixCostTask(fc.FixCostID)

	_, err = s.asynqClient.Enqueue(
		task,
		asynq.ProcessAt(fc.NextRunDate),
		asynq.TaskID("fixcost:"+fc.FixCostID),
	)

	if err != nil {
		log.Printf("[FIX COST] enqueue error: %v", err)
	}
	return &fc, nil
}

func (s *service) Update(ctx context.Context, input *m.FixCostUpdateInput) error {

	exists, err := s.repo.FindByID(ctx, input.FixCostID, input.UserID)
	if err != nil {
		return err
	}

	if exists == nil {
		return gorm.ErrRecordNotFound
	}

	if input.Title != nil {
		exists.Title = *input.Title
	}
	if input.CategoryID != nil {
		exists.CategoryID = *input.CategoryID
	}
	if input.Price != nil {
		exists.Price = *input.Price
	}
	if input.StartDate != nil {
		exists.StartDate = *input.StartDate
	}

	// always update these 2 fields since they affect the schedule
	exists.EndDate = input.EndDate
	exists.RemainingRuns = input.RemainingRuns

	if input.IntervalType != nil {
		exists.IntervalType = e.IntervalType(*input.IntervalType)
	}
	if input.IntervalValue != nil {
		exists.IntervalValue = *input.IntervalValue
	}

	// check schedule change
	scheduleChanged :=
		string(exists.IntervalType) != *input.IntervalType ||
			exists.IntervalValue != *input.IntervalValue ||
			!exists.EndDate.Equal(*input.EndDate)

	err = s.repo.Update(ctx, exists)
	if err != nil {
		return err
	}

	fc, err := s.repo.FindByID(ctx, input.FixCostID, input.UserID)
	if err != nil {
		return err
	}

	if scheduleChanged {

		task, _ := tasks.NewRunFixCostTask(fc.FixCostID)

		_, err = s.asynqClient.Enqueue(
			task,
			asynq.ProcessAt(fc.NextRunDate),
			asynq.TaskID("fixcost:"+fc.FixCostID),
		)

		if err != nil {
			log.Printf("[FIX COST] enqueue error: %v", err)
		}
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id string, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}

func (s *service) RecoverFixCostJobs() {
	ctx := context.Background()

	fixCosts, err := s.repo.FindAllActive(ctx)
	if err != nil {
		log.Printf("[FIX COST SERVICE] recover error: %v", err)
		return
	}

	for _, fc := range fixCosts {

		task, _ := tasks.NewRunFixCostTask(fc.FixCostID)

		_, err := s.asynqClient.Enqueue(
			task,
			asynq.ProcessAt(fc.NextRunDate),
			asynq.TaskID("fixcost:"+fc.FixCostID),
		)

		if err != nil {
			log.Printf("[FIX COST SERVICE] enqueue error: %v", err)
		}
	}
}

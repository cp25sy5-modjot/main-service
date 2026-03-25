package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
	summaryrepo "github.com/cp25sy5-modjot/main-service/internal/summary/repository"
)

type Service interface {
	GetExpenseSummary(ctx context.Context, userID string, period Period) (m.ExpenseSummaryRes, error)
	GetCategorySummary(ctx context.Context, userID string, period Period, date *time.Time) (m.CategorySummaryRes, error)
}

type service struct {
	repo *summaryrepo.Repository
}

func NewService(repo *summaryrepo.Repository) Service {
	return &service{repo}
}

func (s *service) GetExpenseSummary(
	ctx context.Context,
	userID string,
	period Period,
) (m.ExpenseSummaryRes, error) {

	now := time.Now().UTC()

	var (
		format string
		start  time.Time
		end    time.Time
	)

	switch period {

	case Week:

		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		start = time.Date(
			now.Year(),
			now.Month(),
			now.Day()-weekday+1,
			0, 0, 0, 0,
			time.UTC,
		)

		end = start.AddDate(0, 0, 7)

		format = "YYYY-MM-DD"

	case Month:

		// ทุกเดือนของปีนี้
		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(1, 0, 0)

		format = "MM"

	case Year:

		// 3 ปีล่าสุด
		start = time.Date(now.Year()-2, 1, 1, 0, 0, 0, 0, time.UTC)
		end = time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)

		format = "YYYY"

	default:
		return m.ExpenseSummaryRes{}, errors.New("invalid period")
	}

	data, err := s.repo.ExpenseSummary(ctx, userID, format, start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return m.ExpenseSummaryRes{}, err
	}

	data = fillZero(period, start, end, data)

	return m.ExpenseSummaryRes{
		Period: string(period),
		Data:   data,
	}, nil
}

func (s *service) GetCategorySummary(
	ctx context.Context,
	userID string,
	period Period,
	date *time.Time,
) (m.CategorySummaryRes, error) {

	ref := time.Now().UTC()
	if date != nil {
		ref = date.UTC()
	}

	start, end, units, err := resolvePeriodRange(period, ref)
	if err != nil {
		return m.CategorySummaryRes{}, err
	}

	data, err := s.repo.CategorySummary(ctx, userID, start, end)
	if err != nil {
		return m.CategorySummaryRes{}, fmt.Errorf("get category summary: %w", err)
	}

	var total float64
	for _, d := range data {
		total += d.Total
	}

	var avg float64
	if units > 0 {
		avg = total / float64(units)
	}

	return m.CategorySummaryRes{
		Period:  string(period),
		Total:   total,
		Average: avg,
		Data:    data,
	}, nil
}

func resolvePeriodRange(period Period, ref time.Time) (
	start time.Time,
	end time.Time,
	units int,
	err error,
) {

	switch period {

	case Day:
		start = startOfDay(ref)
		end = start.AddDate(0, 0, 1)
		units = 1

	case Week:
		start = startOfWeek(ref)
		end = start.AddDate(0, 0, 7)
		units = 7

	case Month:
		start = startOfMonth(ref)
		end = start.AddDate(0, 1, 0)
		units = int(end.Sub(start).Hours() / 24)

	case Year:
		start = startOfYear(ref)
		end = start.AddDate(1, 0, 0)
		units = 12

	case PastYear:
		const years = 3
		start = startOfYear(ref.AddDate(-years+1, 0, 0))
		end = startOfYear(ref.AddDate(1, 0, 0))
		units = years

	default:
		err = fmt.Errorf("invalid period")
	}

	return
}
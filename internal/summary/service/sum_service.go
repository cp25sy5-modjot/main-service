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
	GetExpenseSummary(ctx context.Context, userID string, period Period, date *time.Time) (m.ExpenseSummaryRes, error)
	GetCategorySummary(ctx context.Context, userID string, period Period, date *time.Time) (m.CategorySummaryRes, error)
}

type service struct {
	repo summaryrepo.Repository
}

func NewService(repo summaryrepo.Repository) Service {
	return &service{repo}
}

func (s *service) GetExpenseSummary(
	ctx context.Context,
	userID string,
	period Period,
	date *time.Time,
) (m.ExpenseSummaryRes, error) {

	start, end, format, err := resolveExpensePeriodRange(period, date.In(time.Local))
	if err != nil {
		return m.ExpenseSummaryRes{}, err
	}

	data, err := s.repo.ExpenseSummary(ctx, userID, format, start.UTC(), end.UTC())
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

	start, end, units, err := resolveCategoryPeriodRange(period, date.In(time.Local))
	if err != nil {
		return m.CategorySummaryRes{}, err
	}

	data, err := s.repo.CategorySummary(ctx, userID, start.UTC(), end.UTC())
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

func resolveExpensePeriodRange(period Period, ref time.Time) (
	start, end time.Time,
	format string,
	err error,
) {

	switch period {

	case Week:

		start = startOfWeek(ref)
		end = start.AddDate(0, 0, 7)
		format = "YYYY-MM-DD"

	case Year:

		start = startOfMonth(ref)
		end = start.AddDate(1, 0, 0)

		format = "MM"

	case Last3Year:

		// 3 ปีล่าสุด
		start = startOfLastNYear(ref, 2)
		end = startOfYear(ref.AddDate(1, 0, 0))

		format = "YYYY"

	default:
		err = errors.New("invalid period")
	}

	return
}

func resolveCategoryPeriodRange(period Period, ref time.Time) (
	start, end time.Time,
	units int,
	err error,
) {

	switch period {
	// single: 1 day, 1 month
	case Day:
		start = startOfDay(ref)
		end = start.AddDate(0, 0, 1)
		units = 1

	case Month:
		start = startOfMonth(ref)
		end = start.AddDate(0, 1, 0)
		units = int(end.Sub(start).Hours() / 24)

	// multiple: 1 week (7 days), 1 year (12 months), and last 3 year (3 years)
	case Week:

		start = startOfWeek(ref)
		end = start.AddDate(0, 0, 7)
		units = 7

	case Year:
		start = startOfYear(ref)
		end = start.AddDate(1, 0, 0)
		units = 12

	case Last3Year:
		const years = 2
		start = startOfLastNYear(ref, years)
		end = startOfYear(ref.AddDate(1, 0, 0))
		units = years + 1

	default:
		err = fmt.Errorf("invalid period")
	}

	return
}

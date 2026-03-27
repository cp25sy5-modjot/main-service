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
	date   *time.Time,
) (m.ExpenseSummaryRes, error) {
loc := date.Location()
	var (
		format string
		start  time.Time
		end    time.Time
	)

	switch period {

	case Week:

		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		start = time.Date(
			date.Year(),
			date.Month(),
			date.Day()-weekday+1,
			0, 0, 0, 0,
			loc,
		)

		end = start.AddDate(0, 0, 7)

		format = "YYYY-MM-DD"

	case Month:

		start = time.Date(date.Year(), 1, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(1, 0, 0)

		format = "MM"

	case Year:

		// 3 ปีล่าสุด
		start = time.Date(date.Year()-2, 1, 1, 0, 0, 0, 0, loc)
		end = time.Date(date.Year()+1, 1, 1, 0, 0, 0, 0, loc)

		format = "YYYY"

	default:
		return m.ExpenseSummaryRes{}, errors.New("invalid period")
	}

	startUTC := start.UTC()
	endUTC := end.UTC()

	data, err := s.repo.ExpenseSummary(ctx, userID, format, startUTC, endUTC)
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

	loc := time.Local // ❗ หรือควร inject จาก client (ดีที่สุด)

	ref := time.Now().In(loc)
	if date != nil {
		ref = date.In(loc)
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
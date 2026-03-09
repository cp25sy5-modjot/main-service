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
	var start time.Time
	var end time.Time
	var units int

	ref := time.Now().UTC()

	if date != nil {
		ref = date.UTC()
	}

	switch period {

	case Week:

		if date != nil {

			start = time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1)

			units = 1

		} else {

			weekday := int(ref.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			start = ref.AddDate(0, 0, -weekday+1)
			end = start.AddDate(0, 0, 7)

			units = 7
		}

	case Month:

		if date != nil {

			start = time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1)

			units = 1

		} else {

			start = time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 1, 0)

			units = 1
		}

	case Year:

		if date != nil {

			start = time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
			end = start.AddDate(0, 0, 1)

			units = 1

		} else {

			start = time.Date(ref.Year()-2, 1, 1, 0, 0, 0, 0, time.UTC)
			end = time.Date(ref.Year()+1, 1, 1, 0, 0, 0, 0, time.UTC)

			units = end.Year() - start.Year()
		}

	default:
		return m.CategorySummaryRes{}, fmt.Errorf("invalid period")
	}

	data, err := s.repo.CategorySummary(ctx, userID, start, end)
	if err != nil {
		return m.CategorySummaryRes{}, fmt.Errorf("get category summary: %w", err)
	}

	var total float64

	for _, d := range data {
		total += d.Total
	}

	avg := 0.0
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

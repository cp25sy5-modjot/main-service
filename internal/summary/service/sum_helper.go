package service

import (
	"fmt"
	"strconv"
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

type Period string

const (
	// for 1 day 1 month
	Day   Period = "day"
	Month Period = "month"

	// for 1 week (7 days), 1 year (12 months), and last 3 year (3 years)
	Week      Period = "week"
	Year      Period = "year"
	Last3Year Period = "last_3_year"
)

func fillZero(period Period, start, end time.Time, data []m.ExpenseSummary) []m.ExpenseSummary {

	resultMap := make(map[string]float64)

	for _, d := range data {
		resultMap[d.Key] = d.Total
	}

	var result []m.ExpenseSummary

	switch period {

	case Week:

		for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {

			key := d.Format("2006-01-02")
			label := d.Weekday().String()[:3]

			result = append(result, m.ExpenseSummary{
				Key:   key,
				Label: label,
				Total: resultMap[key],
			})
		}

	case Year:

		for i := 1; i <= 12; i++ {

			key := fmt.Sprintf("%02d", i)
			label := time.Month(i).String()[:3]

			result = append(result, m.ExpenseSummary{
				Key:   key,
				Label: label,
				Total: resultMap[key],
			})
		}

	case Last3Year:

		for y := start.Year(); y < end.Year(); y++ {

			key := strconv.Itoa(y)

			result = append(result, m.ExpenseSummary{
				Key:   key,
				Label: key,
				Total: resultMap[key],
			})
		}
	}

	return result
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day()-weekday+1,
		0, 0, 0, 0,
		t.Location(),
	)
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func startOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

func startOfLastNYear(t time.Time, n int) time.Time {
	return time.Date(t.Year()-n, 1, 1, 0, 0, 0, 0, t.Location())
}

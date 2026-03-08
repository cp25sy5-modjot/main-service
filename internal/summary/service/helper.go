package service

import (
	"fmt"
	"time"

	m "github.com/cp25sy5-modjot/main-service/internal/domain/model"
)

type Period string

const (
	Week  Period = "week"
	Month Period = "month"
	Year  Period = "year"
)

func fillZero(period Period, start, end time.Time, data []m.ExpenseSummary) []m.ExpenseSummary {

	resultMap := map[string]float64{}

	for _, d := range data {
		resultMap[d.Label] = d.Total
	}

	var result []m.ExpenseSummary

	switch period {

	case Week:

		for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {

			label := d.Format("2006-01-02")

			result = append(result, m.ExpenseSummary{
				Label: label,
				Total: resultMap[label],
			})
		}

	case Month:

		for i := 1; i <= 12; i++ {

			label := fmt.Sprintf("%02d", i)

			result = append(result, m.ExpenseSummary{
				Label: label,
				Total: resultMap[label],
			})
		}

	case Year:

		for y := start.Year(); y <= end.Year(); y++ {

			label := fmt.Sprintf("%d", y)

			result = append(result, m.ExpenseSummary{
				Label: label,
				Total: resultMap[label],
			})
		}

	}

	return result
}

func (p Period) IsValid() bool {
	switch p {
	case Week, Month, Year:
		return true
	default:
		return false
	}
}

func ParsePeriod(s string) (Period, error) {
	if s == "" {
		return Week, nil
	}

	p := Period(s)

	switch p {
	case Week, Month, Year:
		return p, nil
	default:
		return "", fmt.Errorf("invalid period")
	}
}

package service

import (
	"fmt"
	"strconv"
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

	case Month:

		for i := 1; i <= 12; i++ {

			key := fmt.Sprintf("%02d", i)
			label := time.Month(i).String()[:3]

			result = append(result, m.ExpenseSummary{
				Key:   key,
				Label: label,
				Total: resultMap[key],
			})
		}

	case Year:

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

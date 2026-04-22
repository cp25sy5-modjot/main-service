package fixcostsvc

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

func CalculateNextRun(fc e.FixCost) time.Time {
	switch fc.IntervalType {

	case "daily":
		return fc.StartDate.AddDate(0, 0, fc.RunCount*fc.IntervalValue)

	case "weekly":
		return fc.StartDate.AddDate(0, 0, 7*fc.RunCount*fc.IntervalValue)

	case "monthly":
		return calculateMonthly(fc)

	case "yearly":
		return calculateYearly(fc)

	default:
		return fc.NextRunDate
	}
}

func calculateMonthly(fc e.FixCost) time.Time {
	targetMonth := fc.RunCount * fc.IntervalValue

	base := fc.StartDate
	year := base.Year()
	month := int(base.Month()) + targetMonth

	// normalize year/month
	year += (month - 1) / 12
	month = (month-1)%12 + 1

	day := base.Day()

	// หาวันสุดท้ายของเดือน
	lastDay := lastDayOfMonth(year, time.Month(month))
	if day > lastDay {
		day = lastDay
	}

	return time.Date(
		year,
		time.Month(month),
		day,
		base.Hour(),
		base.Minute(),
		base.Second(),
		base.Nanosecond(),
		time.UTC,
	)
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func calculateYearly(fc e.FixCost) time.Time {
	base := fc.StartDate

	targetYear := base.Year() + (fc.RunCount * fc.IntervalValue)
	month := base.Month()
	day := base.Day()

	lastDay := lastDayOfMonth(targetYear, month)
	if day > lastDay {
		day = lastDay
	}

	return time.Date(
		targetYear,
		month,
		day,
		base.Hour(),
		base.Minute(),
		base.Second(),
		base.Nanosecond(),
		time.UTC,
	)
}

func calculateStatus(endDate *time.Time, maxRun *int) e.FixCostStatus {
	now := time.Now()

	// 1. ถ้ามี EndDate และเลยแล้ว → หมดอายุ
	if endDate != nil && now.After(*endDate) {
		return e.FixCostStatusFinished
	}

	// 2. ถ้ามี MaxRun และเหลือ <= 0 → หมด
	if maxRun != nil && *maxRun <= 0 {
		return e.FixCostStatusFinished
	}

	// 3. default → active
	return e.FixCostStatusActive
}

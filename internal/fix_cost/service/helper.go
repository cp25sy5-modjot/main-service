package fixcostsvc

import (
	"time"

	e "github.com/cp25sy5-modjot/main-service/internal/domain/entity"
)

func CalculateNextRun(fc e.FixCost) time.Time {
	switch fc.IntervalType {

	case "daily":
		return fc.NextRunDate.AddDate(0, 0, fc.IntervalValue)

	case "weekly":
		return fc.NextRunDate.AddDate(0, 0, 7*fc.IntervalValue)

	case "monthly":
		return fc.NextRunDate.AddDate(0, fc.IntervalValue, 0)

	case "yearly":
		return fc.NextRunDate.AddDate(fc.IntervalValue, 0, 0)

	default:
		return fc.NextRunDate
	}
}

func calculateStatus(endDate *time.Time, remainingRuns *int) e.FixCostStatus {
	now := time.Now()

	// 1. ถ้ามี EndDate และเลยแล้ว → หมดอายุ
	if endDate != nil && now.After(*endDate) {
		return e.FixCostStatusFinished
	}

	// 2. ถ้ามี RemainingRuns และเหลือ <= 0 → หมด
	if remainingRuns != nil && *remainingRuns <= 0 {
		return e.FixCostStatusFinished
	}

	// 3. default → active
	return e.FixCostStatusActive
}

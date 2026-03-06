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

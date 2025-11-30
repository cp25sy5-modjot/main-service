package utils

import "time"

func GetStartAndEndOfMonth(t time.Time) (time.Time, time.Time) {
	loc := t.Location()

	// First day of current month at 00:00
	firstOfCurrent := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)

	// Last day of previous month = firstOfCurrent - 1 day
	lastOfPrevious := firstOfCurrent.AddDate(0, 0, -1)

	// Last day of current month = first of next month - 1 day
	firstOfNext := firstOfCurrent.AddDate(0, 1, 0)
	lastOfCurrent := firstOfNext.AddDate(0, 0, -1)

	// Set final times to 17:00:00
	startRange := time.Date(
		lastOfPrevious.Year(), lastOfPrevious.Month(), lastOfPrevious.Day(),
		17, 0, 0, 0,
		loc,
	)

	endRange := time.Date(
		lastOfCurrent.Year(), lastOfCurrent.Month(), lastOfCurrent.Day(),
		17, 0, 0, 0,
		loc,
	)

	return startRange, endRange
}

package utils

import "time"

// func GetStartAndEndOfMonth(t time.Time) (time.Time, time.Time) {
// 	loc := t.Location()

// 	// First day of current month at 00:00
// 	firstOfCurrent := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)

// 	// Last day of previous month = firstOfCurrent - 1 day
// 	lastOfPrevious := firstOfCurrent.AddDate(0, 0, -1)

// 	// Last day of current month = first of next month - 1 day
// 	firstOfNext := firstOfCurrent.AddDate(0, 1, 0)
// 	lastOfCurrent := firstOfNext.AddDate(0, 0, -1)

// 	// Set final times to 17:00:00
// 	startRange := time.Date(
// 		lastOfPrevious.Year(), lastOfPrevious.Month(), lastOfPrevious.Day(),
// 		17, 0, 0, 0,
// 		loc,
// 	)

// 	endRange := time.Date(
// 		lastOfCurrent.Year(), lastOfCurrent.Month(), lastOfCurrent.Day(),
// 		17, 0, 0, 0,
// 		loc,
// 	)

// 	return startRange, endRange
// }

// func GetStartAndEndOfPreviousMonth(t time.Time) (time.Time, time.Time) {
// 	loc := t.Location()

// 	// First day of current month at 00:00
// 	firstOfCurrent := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)

// 	// Last day of previous month = firstOfCurrent - 1 day
// 	lastOfPrevious := firstOfCurrent.AddDate(0, 0, -1)

// 	// First day of previous month
// 	firstOfPrevious := time.Date(lastOfPrevious.Year(), lastOfPrevious.Month(), 1, 0, 0, 0, 0, loc)

// 	// Set final times to 17:00:00
// 	startRange := time.Date(
// 		firstOfPrevious.Year(), firstOfPrevious.Month(), firstOfPrevious.Day(),
// 		17, 0, 0, 0,
// 		loc,
// 	)

// 	endRange := time.Date(
// 		lastOfPrevious.Year(), lastOfPrevious.Month(), lastOfPrevious.Day(),
// 		17, 0, 0, 0,
// 		loc,
// 	)

// 	return startRange, endRange
// }

func GetStartAndEndOfMonth(t time.Time) (time.Time, time.Time) {
	loc := t.Location()

	start := time.Date(
		t.Year(), t.Month(), 1,
		0, 0, 0, 0,
		loc,
	)

	end := start.AddDate(0, 1, 0)

	return start, end
}

func GetStartAndEndOfPreviousMonth(t time.Time) (time.Time, time.Time) {
	loc := t.Location()

	// First day of current month
	firstOfCurrent := time.Date(
		t.Year(), t.Month(), 1,
		0, 0, 0, 0,
		loc,
	)

	start := firstOfCurrent.AddDate(0, -1, 0)
	end := firstOfCurrent

	return start, end
}

// NormalizeToUTC ensures time is stored in UTC
func NormalizeToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

func ParseRFC3339ToUTC(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func ParseDateInLocationToUTC(date string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

func GetMonthRangeUTC(t time.Time) (time.Time, time.Time) {
	start := time.Date(
		t.Year(), t.Month(), 1,
		0, 0, 0, 0,
		time.UTC,
	)
	end := start.AddDate(0, 1, 0)
	return start, end
}

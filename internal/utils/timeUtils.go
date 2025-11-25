package utils

import "time"

const DefaultTZ = "Asia/Bangkok"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func LoadLocationOrDefault(tz string) *time.Location {
	if tz == "" {
		tz = DefaultTZ
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc, _ = time.LoadLocation(DefaultTZ)
	}
	return loc
}

func ToUserLocal(t time.Time, tz string) time.Time {
	loc := LoadLocationOrDefault(tz)
	return t.In(loc)
}

func NormalizeToUTC(t time.Time, tz string) time.Time {
	if t.Location() == time.UTC || t.Location().String() == "Local" {
		loc := LoadLocationOrDefault(tz)
		// Treat the time as if it's user local time
		t = time.Date(
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(),
			t.Nanosecond(),
			loc,
		)
	}
	return t.UTC()
}

package task

import "time"

func floorDate(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func mergeDateAndClock(date time.Time, clockSource time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		clockSource.Hour(),
		clockSource.Minute(),
		clockSource.Second(),
		clockSource.Nanosecond(),
		time.UTC,
	)
}

func minTime(a time.Time, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func maxTime(a time.Time, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

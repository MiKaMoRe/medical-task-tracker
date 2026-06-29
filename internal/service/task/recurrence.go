package task

import (
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func hasRecurringRule(t *domaintask.Task) bool {
	return t != nil && t.IsRecurring && t.Recurring != nil && t.Recurring.Rule != nil
}

func recurringRangeForPeriod(t *domaintask.Task, from time.Time, to time.Time) (time.Time, time.Time, bool) {
	rangeEnd := to
	if t.Recurring.EndDate != nil {
		recurringEnd := time.Time(*t.Recurring.EndDate).UTC()
		rangeEnd = minTime(recurringEnd, to)
	}
	rangeStart := maxTime(time.Time(t.Date).UTC(), from)
	if rangeStart.After(rangeEnd) {
		return time.Time{}, time.Time{}, false
	}
	return rangeStart, rangeEnd, true
}

func buildRecurringOccurrencesForPeriod(t *domaintask.Task, from time.Time, to time.Time) []*domaintask.Task {
	rangeStart, rangeEnd, ok := recurringRangeForPeriod(t, from, to)
	if !ok {
		return nil
	}

	result := make([]*domaintask.Task, 0)
	for d := floorDate(rangeStart); !d.After(floorDate(rangeEnd)); d = d.AddDate(0, 0, 1) {
		occurrence := mergeDateAndClock(d, time.Time(t.Date).UTC())
		if !matchesRecurringRule(occurrence, t) {
			continue
		}

		doneDateKey := floorDate(occurrence).Format(time.DateOnly)
		_, isDone := t.DoneDates[doneDateKey]
		result = append(result, cloneTaskWithDate(t, occurrence, isDone))
	}
	return result
}

func matchesRecurringRule(date time.Time, t *domaintask.Task) bool {
	if !hasRecurringRule(t) {
		return false
	}
	return t.Recurring.Rule.Matches(date.UTC(), time.Time(t.Date).UTC())
}

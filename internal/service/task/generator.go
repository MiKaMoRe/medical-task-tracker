package task

import (
	"sort"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func generateTasksForPeriod(tasks []*domaintask.Task, from time.Time, to time.Time) []*domaintask.Task {
	result := make([]*domaintask.Task, 0, len(tasks))
	for _, t := range tasks {
		if t == nil {
			continue
		}

		if !t.IsRecurring || t.Recurring == nil || t.Recurring.Rule == nil {
			taskDate := time.Time(t.Date).UTC()
			if !taskDate.Before(from) && !taskDate.After(to) {
				result = append(result, cloneTaskWithDate(t, taskDate, t.IsDone))
			}
			continue
		}

		rangeEnd := to
		if t.Recurring.EndDate != nil {
			recurringEnd := time.Time(*t.Recurring.EndDate).UTC()
			rangeEnd = minTime(recurringEnd, to)
		}
		rangeStart := maxTime(time.Time(t.Date).UTC(), from)
		if rangeStart.After(rangeEnd) {
			continue
		}

		for d := floorDate(rangeStart); !d.After(floorDate(rangeEnd)); d = d.AddDate(0, 0, 1) {
			occurrence := mergeDateAndClock(d, time.Time(t.Date).UTC())
			if matchesRecurringRule(occurrence, t) {
				doneDateKey := floorDate(occurrence).Format(time.DateOnly)
				_, isDone := t.DoneDates[doneDateKey]
				result = append(result, cloneTaskWithDate(t, occurrence, isDone))
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		left := time.Time(result[i].Date)
		right := time.Time(result[j].Date)
		if left.Equal(right) {
			return result[i].ID.Int() < result[j].ID.Int()
		}
		return left.Before(right)
	})

	return result
}

func matchesRecurringRule(date time.Time, t *domaintask.Task) bool {
	switch rule := t.Recurring.Rule.(type) {
	case domaintask.WeeklyRecurrenceRule:
		return int(date.Weekday()) == rule.WeekDay
	case domaintask.MonthlyRecurrenceRule:
		return date.Day() == rule.MonthDay
	case domaintask.YearlyRecurrenceRule:
		return int(date.Month()) == rule.Month && date.Day() == rule.Day
	case domaintask.BiweeklyRecurrenceRule:
		year, week := date.ISOWeek()
		_ = year
		isOddWeek := week%2 == 1
		return int(date.Weekday()) == rule.WeekDay && isOddWeek == rule.IsOdd
	case domaintask.ShiftRecurrenceRule:
		base := floorDate(time.Time(t.Date).UTC())
		current := floorDate(date.UTC())
		daysFromStart := int(current.Sub(base).Hours() / 24)
		cycle := rule.NumberOfTaskDays + rule.NumberOfShiftDays
		if cycle <= 0 {
			return false
		}
		dayInCycle := daysFromStart % cycle
		return dayInCycle >= 0 && dayInCycle < rule.NumberOfTaskDays
	case domaintask.ParityRecurrenceRule:
		isEvenDay := date.Day()%2 == 0
		return isEvenDay == rule.IsEven
	default:
		return false
	}
}

func cloneTaskWithDate(t *domaintask.Task, date time.Time, isDone bool) *domaintask.Task {
	cloned := *t
	cloned.Date = domaintask.Date(date.UTC())
	cloned.IsDone = isDone
	return &cloned
}

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

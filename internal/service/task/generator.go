package task

import (
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func generateTasksForPeriod(tasks []*domaintask.Task, from time.Time, to time.Time) []*domaintask.Task {
	result := make([]*domaintask.Task, 0, len(tasks))
	for _, t := range tasks {
		if t == nil {
			continue
		}

		if !hasRecurringRule(t) {
			taskDate := time.Time(t.Date).UTC()
			if !taskDate.Before(from) && !taskDate.After(to) {
				result = append(result, cloneTaskWithDate(t, taskDate, t.IsDone))
			}
			continue
		}

		result = append(result, buildRecurringOccurrencesForPeriod(t, from, to)...)
	}

	sortTasksByDateThenID(result)

	return result
}

func cloneTaskWithDate(t *domaintask.Task, date time.Time, isDone bool) *domaintask.Task {
	cloned := *t
	cloned.Date = domaintask.Date(date.UTC())
	cloned.IsDone = isDone
	return &cloned
}

package task

import (
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func assignTaskStatuses(tasks []*domaintask.Task, now time.Time) {
	for _, t := range tasks {
		if t == nil {
			continue
		}
		isDone := t.IsDone
		if t.IsRecurring && len(t.DoneDates) > 0 {
			_, isDone = t.DoneDates[floorDate(time.Time(t.Date).UTC()).Format(time.DateOnly)]
		}
		t.Status = domaintask.ResolveTaskStatus(time.Time(t.Date).UTC(), isDone, now)
	}
}

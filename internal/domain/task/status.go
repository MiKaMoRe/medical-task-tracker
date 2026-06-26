package task

import "time"

type TaskStatus string

const (
	TaskStatusPlanned TaskStatus = "planned"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusExpired TaskStatus = "expired"
)

func ResolveTaskStatus(taskDate time.Time, isDone bool, now time.Time) TaskStatus {
	if isDone {
		return TaskStatusDone
	}

	taskDay := dateOnlyUTC(taskDate)
	nowDay := dateOnlyUTC(now)
	if taskDay.Before(nowDay) {
		return TaskStatusExpired
	}

	return TaskStatusPlanned
}

func dateOnlyUTC(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

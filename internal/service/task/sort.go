package task

import (
	"sort"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func sortTasksByDateThenID(tasks []*domaintask.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		left := time.Time(tasks[i].Date)
		right := time.Time(tasks[j].Date)
		if left.Equal(right) {
			return tasks[i].ID.Int() < tasks[j].ID.Int()
		}
		return left.Before(right)
	})
}

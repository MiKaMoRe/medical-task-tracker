package task

import (
	"context"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (s *TaskService) GetTask(ctx context.Context, id task.ID) (*task.Task, error) {
	t, err := s.repo.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	assignTaskStatuses([]*task.Task{t}, time.Now().UTC())
	return t, nil
}

func (s *TaskService) GetTasks(ctx context.Context, from time.Time, to time.Time, statuses []task.TaskStatus) ([]*task.Task, error) {
	if from.IsZero() || to.IsZero() {
		return nil, apperrors.NewAppError("from and to are required")
	}
	if from.After(to) {
		return nil, apperrors.NewAppError("from must be less than or equal to to")
	}

	rawTasks, err := s.repo.GetTasksForPeriod(ctx, from.UTC(), to.UTC())
	if err != nil {
		return nil, err
	}

	generated := generateTasksForPeriod(rawTasks, from.UTC(), to.UTC())
	assignTaskStatuses(generated, time.Now().UTC())
	if len(statuses) > 0 {
		generated = filterTasksByStatuses(generated, statuses)
	}
	return generated, nil
}

func filterTasksByStatuses(tasks []*task.Task, statuses []task.TaskStatus) []*task.Task {
	allowed := make(map[task.TaskStatus]struct{}, len(statuses))
	for _, status := range statuses {
		allowed[status] = struct{}{}
	}

	filtered := make([]*task.Task, 0, len(tasks))
	for _, t := range tasks {
		if t == nil {
			continue
		}
		if _, ok := allowed[t.Status]; ok {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

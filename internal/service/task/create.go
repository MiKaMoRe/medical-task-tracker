package task

import (
	"context"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (s *TaskService) CreateTask(ctx context.Context, task *task.Task) (*task.Task, error) {
	return s.repo.CreateTask(ctx, task)
}

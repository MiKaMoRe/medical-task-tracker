package task

import (
	"context"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	GetTask(ctx context.Context, id string) (*task.Task, error)
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

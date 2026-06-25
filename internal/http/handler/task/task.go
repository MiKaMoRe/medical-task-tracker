package task

import (
	"context"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	GetTask(ctx context.Context, id int) (*task.Task, error)
}

type TaskHandler struct {
	srvc   TaskService
	logger logger.Logger
}

func NewTaskHandler(srvc TaskService, logger logger.Logger) *TaskHandler {
	return &TaskHandler{
		srvc:   srvc,
		logger: logger,
	}
}

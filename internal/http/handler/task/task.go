package task

import (
	"context"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/logger"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	UpdateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	DeleteTask(ctx context.Context, id task.ID) error
	AddTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	RemoveTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	GetTask(ctx context.Context, id task.ID) (*task.Task, error)
	GetTasks(ctx context.Context, from time.Time, to time.Time, statuses []task.TaskStatus) ([]*task.Task, error)
	MarkTaskDone(ctx context.Context, id task.ID, occurrenceDate *time.Time) error
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

package task

import (
	"context"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/db/txrunner"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	UpdateTask(ctx context.Context, task *task.Task) (*task.Task, error)
	DeleteTask(ctx context.Context, id task.ID) error
	AddTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	RemoveTaskTags(ctx context.Context, id task.ID, tags []string) (*task.Task, error)
	GetTask(ctx context.Context, id task.ID) (*task.Task, error)
	GetTasksForPeriod(ctx context.Context, from time.Time, to time.Time) ([]*task.Task, error)
	MarkTaskDone(ctx context.Context, id task.ID) error
	MarkTaskOccurrenceDone(ctx context.Context, id task.ID, occurrenceDate time.Time) error
}

type TaskService struct {
	repo TaskRepository
	tx   txrunner.Runner
}

func NewTaskService(repo TaskRepository, tx txrunner.Runner) *TaskService {
	return &TaskService{repo: repo, tx: tx}
}

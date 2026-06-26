package task

import (
	"context"
	"errors"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (s *TaskService) CreateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	vm := validateTask(t)
	if vm.Err() != nil {
		return nil, vm.Err()
	}

	var created *task.Task
	err := s.withWriteTx(ctx, func(ctx context.Context) error {
		var err error
		created, err = s.repo.CreateTask(ctx, t)
		return err
	})
	if err != nil {
		return nil, err
	}

	created.Status = task.TaskStatusPlanned
	return created, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, t *task.Task) (*task.Task, error) {
	vm := validateTask(t)
	if vm.Err() != nil {
		return nil, vm.Err()
	}

	var updated *task.Task
	err := s.withWriteTx(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.UpdateTask(ctx, t)
		return err
	})
	if err != nil {
		return nil, err
	}

	assignTaskStatuses([]*task.Task{updated}, time.Now().UTC())
	return updated, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id task.ID) error {
	return s.withWriteTx(ctx, func(ctx context.Context) error {
		return s.repo.DeleteTask(ctx, id)
	})
}

func validateTask(t *task.Task) *apperrors.ValidationMap {
	vm := apperrors.NewValidationMap()
	if t.IsRecurring && t.Recurring == nil {
		vm.Add("task.recurring", errors.New("recurring task is required when is_recurring is true"))
	}
	if t.Recurring != nil && t.Recurring.Rule == nil {
		vm.Add("task.recurring.rule", errors.New("recurring rule is required"))
	}

	return vm
}

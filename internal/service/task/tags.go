package task

import (
	"context"
	"errors"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	domaintag "github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (s *TaskService) AddTaskTags(ctx context.Context, id domaintask.ID, tags []string) (*domaintask.Task, error) {
	if err := validateTagNames(tags); err != nil {
		return nil, err
	}

	var updated *domaintask.Task
	err := s.withWriteTx(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.AddTaskTags(ctx, id, tags)
		return err
	})
	if err != nil {
		return nil, err
	}

	assignTaskStatuses([]*domaintask.Task{updated}, time.Now().UTC())
	return updated, nil
}

func (s *TaskService) RemoveTaskTags(ctx context.Context, id domaintask.ID, tags []string) (*domaintask.Task, error) {
	if err := validateTagNames(tags); err != nil {
		return nil, err
	}

	var updated *domaintask.Task
	err := s.withWriteTx(ctx, func(ctx context.Context) error {
		var err error
		updated, err = s.repo.RemoveTaskTags(ctx, id, tags)
		return err
	})
	if err != nil {
		return nil, err
	}

	assignTaskStatuses([]*domaintask.Task{updated}, time.Now().UTC())
	return updated, nil
}

func validateTagNames(names []string) error {
	vm := apperrors.NewValidationMap()
	if len(names) == 0 {
		vm.Add("task.tags", errors.New("at least one tag is required"))
		return vm.Err()
	}
	_, errs := domaintag.NewTags(names)
	vm.Add("task.tags", errs...)
	return vm.Err()
}

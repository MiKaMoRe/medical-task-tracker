package task

import (
	"context"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (s *TaskService) MarkTaskDone(ctx context.Context, id domaintask.ID, occurrenceDate *time.Time) error {
	return s.withWriteTx(ctx, func(ctx context.Context) error {
		targetTask, err := s.repo.GetTask(ctx, id)
		if err != nil {
			return err
		}

		if !targetTask.IsRecurring {
			return s.repo.MarkTaskDone(ctx, id)
		}

		if occurrenceDate == nil {
			return apperrors.NewAppError("occurrence_date is required for recurring task")
		}

		occurrenceDay := floorDate(occurrenceDate.UTC())
		if err := validateRecurringOccurrenceDate(targetTask, occurrenceDay); err != nil {
			return err
		}

		return s.repo.MarkTaskOccurrenceDone(ctx, id, occurrenceDay)
	})
}

func validateRecurringOccurrenceDate(t *domaintask.Task, occurrenceDate time.Time) error {
	taskStart := floorDate(time.Time(t.Date).UTC())
	if occurrenceDate.Before(taskStart) {
		return apperrors.NewAppError("occurrence_date cannot be earlier than task start date")
	}
	if !hasRecurringRule(t) {
		return apperrors.NewAppError("recurring task has invalid recurrence rule")
	}

	if t.Recurring.EndDate != nil {
		recurringEnd := floorDate(time.Time(*t.Recurring.EndDate).UTC())
		if occurrenceDate.After(recurringEnd) {
			return apperrors.NewAppError("occurrence_date is after recurrence end date")
		}
	}

	taskDateTime := mergeDateAndClock(occurrenceDate, time.Time(t.Date).UTC())
	if !matchesRecurringRule(taskDateTime, t) {
		return apperrors.NewAppError("occurrence_date does not match recurrence rule")
	}
	return nil
}

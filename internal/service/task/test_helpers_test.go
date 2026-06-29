package task

import (
	"context"
	"strings"
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

type taskRepoMock struct {
	getTaskFn                func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error)
	getTasksForPeriodFn      func(ctx context.Context, from time.Time, to time.Time) ([]*domaintask.Task, error)
	markTaskDoneFn           func(ctx context.Context, id domaintask.ID) error
	markTaskOccurrenceDoneFn func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error
}

func recurringEndDate(t time.Time) *domaintask.Date {
	date := domaintask.Date(t)
	return &date
}

func mustUTC(year int, month time.Month, day int, hour int, minute int, second int) time.Time {
	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}

func assertErrorContains(t *testing.T, err error, message string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", message)
	}
	if got := err.Error(); got == "" || !strings.Contains(got, message) {
		t.Fatalf("expected error containing %q, got: %v", message, err)
	}
}

func (m *taskRepoMock) CreateTask(ctx context.Context, task *domaintask.Task) (*domaintask.Task, error) {
	panic("unexpected CreateTask call")
}

func (m *taskRepoMock) UpdateTask(ctx context.Context, task *domaintask.Task) (*domaintask.Task, error) {
	panic("unexpected UpdateTask call")
}

func (m *taskRepoMock) DeleteTask(ctx context.Context, id domaintask.ID) error {
	panic("unexpected DeleteTask call")
}

func (m *taskRepoMock) AddTaskTags(ctx context.Context, id domaintask.ID, tags []string) (*domaintask.Task, error) {
	panic("unexpected AddTaskTags call")
}

func (m *taskRepoMock) RemoveTaskTags(ctx context.Context, id domaintask.ID, tags []string) (*domaintask.Task, error) {
	panic("unexpected RemoveTaskTags call")
}

func (m *taskRepoMock) GetTask(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
	if m.getTaskFn == nil {
		panic("unexpected GetTask call")
	}
	return m.getTaskFn(ctx, id)
}

func (m *taskRepoMock) GetTasksForPeriod(ctx context.Context, from time.Time, to time.Time) ([]*domaintask.Task, error) {
	if m.getTasksForPeriodFn == nil {
		panic("unexpected GetTasksForPeriod call")
	}
	return m.getTasksForPeriodFn(ctx, from, to)
}

func (m *taskRepoMock) MarkTaskDone(ctx context.Context, id domaintask.ID) error {
	if m.markTaskDoneFn == nil {
		panic("unexpected MarkTaskDone call")
	}
	return m.markTaskDoneFn(ctx, id)
}

func (m *taskRepoMock) MarkTaskOccurrenceDone(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
	if m.markTaskOccurrenceDoneFn == nil {
		panic("unexpected MarkTaskOccurrenceDone call")
	}
	return m.markTaskOccurrenceDoneFn(ctx, id, occurrenceDate)
}

package task

import (
	"context"
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func TestTaskServiceMarkTaskDone_NonRecurringMarksTask(t *testing.T) {
	taskID := domaintask.IDFromInt(99)
	var markTaskDoneCalled bool
	var markOccurrenceCalled bool

	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(mustUTC(2026, time.June, 26, 9, 0, 0)),
				IsRecurring: false,
			}, nil
		},
		markTaskDoneFn: func(ctx context.Context, id domaintask.ID) error {
			markTaskDoneCalled = true
			return nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			markOccurrenceCalled = true
			return nil
		},
	}
	svc := &TaskService{repo: repo}
	if err := svc.MarkTaskDone(context.Background(), taskID, nil); err != nil {
		t.Fatalf("MarkTaskDone returned error: %v", err)
	}
	if !markTaskDoneCalled {
		t.Fatalf("expected MarkTaskDone to be called for non-recurring task")
	}
	if markOccurrenceCalled {
		t.Fatalf("expected MarkTaskOccurrenceDone not to be called for non-recurring task")
	}
}

func TestTaskServiceMarkTaskDone_RecurringMarksOnlyOccurrence(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := mustUTC(2026, time.June, 26, 9, 0, 0)
	end := time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)
	occurrenceInput := time.Date(2026, time.June, 27, 18, 30, 0, 0, time.UTC)
	expectedOccurrenceDay := mustUTC(2026, time.June, 27, 0, 0, 0)

	var markTaskDoneCalled bool
	var markOccurrenceCalled bool
	var markedID domaintask.ID
	var markedDate time.Time

	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(start),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeShift,
					EndDate:       recurringEndDate(end),
					Rule: domaintask.ShiftRecurrenceRule{
						NumberOfTaskDays:  1,
						NumberOfShiftDays: 0,
					},
				},
			}, nil
		},
		markTaskDoneFn: func(ctx context.Context, id domaintask.ID) error {
			markTaskDoneCalled = true
			return nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			markOccurrenceCalled = true
			markedID = id
			markedDate = occurrenceDate
			return nil
		},
	}
	svc := &TaskService{repo: repo}

	if err := svc.MarkTaskDone(context.Background(), taskID, &occurrenceInput); err != nil {
		t.Fatalf("MarkTaskDone returned error: %v", err)
	}

	if markTaskDoneCalled {
		t.Fatalf("expected MarkTaskDone not to be called for recurring task")
	}
	if !markOccurrenceCalled {
		t.Fatalf("expected MarkTaskOccurrenceDone to be called for recurring task")
	}
	if markedID != taskID {
		t.Fatalf("expected marked id %d, got %d", taskID.Int(), markedID.Int())
	}
	if !markedDate.Equal(expectedOccurrenceDay) {
		t.Fatalf("expected marked date %s, got %s", expectedOccurrenceDay.Format(time.RFC3339), markedDate.Format(time.RFC3339))
	}
}

func TestTaskServiceMarkTaskDone_RecurringRequiresOccurrenceDate(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(mustUTC(2026, time.June, 26, 9, 0, 0)),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeShift,
					EndDate:       recurringEndDate(time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)),
					Rule: domaintask.ShiftRecurrenceRule{
						NumberOfTaskDays:  1,
						NumberOfShiftDays: 0,
					},
				},
			}, nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			t.Fatalf("did not expect MarkTaskOccurrenceDone call")
			return nil
		},
	}
	svc := &TaskService{repo: repo}

	err := svc.MarkTaskDone(context.Background(), taskID, nil)
	assertErrorContains(t, err, "occurrence_date is required for recurring task")
}

func TestTaskServiceMarkTaskDone_RecurringRejectsOccurrenceBeforeStart(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := mustUTC(2026, time.June, 26, 9, 0, 0)
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(start),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeShift,
					EndDate:       recurringEndDate(time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)),
					Rule: domaintask.ShiftRecurrenceRule{
						NumberOfTaskDays:  1,
						NumberOfShiftDays: 0,
					},
				},
			}, nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			t.Fatalf("did not expect MarkTaskOccurrenceDone call")
			return nil
		},
	}
	svc := &TaskService{repo: repo}
	occurrence := start.AddDate(0, 0, -1)

	err := svc.MarkTaskDone(context.Background(), taskID, &occurrence)
	assertErrorContains(t, err, "earlier than task start date")
}

func TestTaskServiceMarkTaskDone_RecurringRejectsOccurrenceAfterEnd(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := mustUTC(2026, time.June, 26, 9, 0, 0)
	end := time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(start),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeShift,
					EndDate:       recurringEndDate(end),
					Rule: domaintask.ShiftRecurrenceRule{
						NumberOfTaskDays:  1,
						NumberOfShiftDays: 0,
					},
				},
			}, nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			t.Fatalf("did not expect MarkTaskOccurrenceDone call")
			return nil
		},
	}
	svc := &TaskService{repo: repo}
	occurrence := end.AddDate(0, 0, 1)

	err := svc.MarkTaskDone(context.Background(), taskID, &occurrence)
	assertErrorContains(t, err, "after recurrence end date")
}

func TestTaskServiceMarkTaskDone_RecurringWithoutEndDateAllowsFutureOccurrence(t *testing.T) {
	taskID := domaintask.IDFromInt(220)
	start := mustUTC(2026, time.June, 22, 9, 0, 0)
	occurrence := mustUTC(2028, time.June, 19, 12, 0, 0)

	var markOccurrenceCalled bool
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(start),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeWeekly,
					EndDate:       nil,
					Rule: domaintask.WeeklyRecurrenceRule{
						WeekDay: 1,
					},
				},
			}, nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			markOccurrenceCalled = true
			return nil
		},
	}
	svc := &TaskService{repo: repo}

	if err := svc.MarkTaskDone(context.Background(), taskID, &occurrence); err != nil {
		t.Fatalf("expected no error for recurring task without end date, got: %v", err)
	}
	if !markOccurrenceCalled {
		t.Fatalf("expected MarkTaskOccurrenceDone to be called")
	}
}

func TestTaskServiceMarkTaskDone_RecurringRejectsOccurrenceNotMatchingRule(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := mustUTC(2026, time.June, 22, 9, 0, 0)
	end := time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(start),
				IsRecurring: true,
				Recurring: &domaintask.RecurringTask{
					RecurringType: domaintask.RecurringTypeWeekly,
					EndDate:       recurringEndDate(end),
					Rule: domaintask.WeeklyRecurrenceRule{
						WeekDay: 1,
					},
				},
			}, nil
		},
		markTaskOccurrenceDoneFn: func(ctx context.Context, id domaintask.ID, occurrenceDate time.Time) error {
			t.Fatalf("did not expect MarkTaskOccurrenceDone call")
			return nil
		},
	}
	svc := &TaskService{repo: repo}
	occurrence := mustUTC(2026, time.June, 23, 13, 0, 0)

	err := svc.MarkTaskDone(context.Background(), taskID, &occurrence)
	assertErrorContains(t, err, "does not match recurrence rule")
}

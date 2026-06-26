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

func TestGenerateTasksForPeriod_IncludesFromAndToBoundaries(t *testing.T) {
	from := time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.July, 3, 23, 59, 59, 0, time.UTC)

	tasks := []*domaintask.Task{
		{ID: domaintask.IDFromInt(1), Date: domaintask.Date(from), IsRecurring: false},
		{ID: domaintask.IDFromInt(2), Date: domaintask.Date(to), IsRecurring: false},
		{ID: domaintask.IDFromInt(3), Date: domaintask.Date(from.Add(-time.Second)), IsRecurring: false},
		{ID: domaintask.IDFromInt(4), Date: domaintask.Date(to.Add(time.Second)), IsRecurring: false},
	}

	got := generateTasksForPeriod(tasks, from, to)
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks in range, got %d", len(got))
	}
	if got[0].ID.Int() != 1 || got[1].ID.Int() != 2 {
		t.Fatalf("expected boundary task ids [1,2], got [%d,%d]", got[0].ID.Int(), got[1].ID.Int())
	}
}

func TestGenerateTasksForPeriod_DailyRecurringGeneratesEightWithSameID(t *testing.T) {
	from := time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.July, 3, 23, 59, 59, 999999999, time.UTC)
	start := time.Date(2026, time.June, 26, 14, 4, 19, 0, time.UTC)
	id := domaintask.IDFromInt(22)

	recurring := &domaintask.Task{
		ID:          id,
		Date:        domaintask.Date(start),
		IsRecurring: true,
		Recurring: &domaintask.RecurringTask{
			RecurringType: domaintask.RecurringTypeShift,
			EndDate:       recurringEndDate(time.Date(2026, time.August, 5, 23, 59, 59, 0, time.UTC)),
			Rule: domaintask.ShiftRecurrenceRule{
				NumberOfTaskDays:  1,
				NumberOfShiftDays: 0,
			},
		},
		DoneDates: map[string]struct{}{},
	}

	got := generateTasksForPeriod([]*domaintask.Task{recurring}, from, to)
	if len(got) != 8 {
		t.Fatalf("expected 8 generated tasks, got %d", len(got))
	}

	for i := range got {
		if got[i].ID != id {
			t.Fatalf("expected generated task id %d, got %d", id.Int(), got[i].ID.Int())
		}

		expectedDate := start.AddDate(0, 0, i)
		if !time.Time(got[i].Date).Equal(expectedDate) {
			t.Fatalf("expected date %s at index %d, got %s", expectedDate.Format(time.RFC3339), i, time.Time(got[i].Date).Format(time.RFC3339))
		}
	}
}

func TestGenerateTasksForPeriod_RecurringWithoutEndDateUsesQueryUpperBound(t *testing.T) {
	from := time.Date(2026, time.June, 26, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)
	start := time.Date(2026, time.June, 26, 9, 0, 0, 0, time.UTC)

	recurring := &domaintask.Task{
		ID:          domaintask.IDFromInt(223),
		Date:        domaintask.Date(start),
		IsRecurring: true,
		Recurring: &domaintask.RecurringTask{
			RecurringType: domaintask.RecurringTypeShift,
			EndDate:       nil,
			Rule: domaintask.ShiftRecurrenceRule{
				NumberOfTaskDays:  1,
				NumberOfShiftDays: 0,
			},
		},
		DoneDates: map[string]struct{}{},
	}

	got := generateTasksForPeriod([]*domaintask.Task{recurring}, from, to)
	if len(got) != 5 {
		t.Fatalf("expected 5 generated tasks, got %d", len(got))
	}
	if !time.Time(got[len(got)-1].Date).Equal(time.Date(2026, time.June, 30, 9, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected last generated date to match query upper bound day")
	}
}

func TestGenerateTasksForPeriod_AllRecurringTypes(t *testing.T) {
	from := time.Date(2026, time.June, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)

	tasks := []*domaintask.Task{
		{
			ID:          domaintask.IDFromInt(101),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 9, 0, 0, 0, time.UTC)), // Monday
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeWeekly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.WeeklyRecurrenceRule{
					WeekDay: 1, // Monday
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(102),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 10, 0, 0, 0, time.UTC)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeMonthly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.MonthlyRecurrenceRule{
					MonthDay: 15,
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(103),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 11, 0, 0, 0, time.UTC)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeYearly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.YearlyRecurrenceRule{
					Month: 6,
					Day:   7,
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(104),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC)), // Week 23 (odd)
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeBiweekly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.BiweeklyRecurrenceRule{
					IsOdd:   true,
					WeekDay: 1, // Monday
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(105),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 13, 0, 0, 0, time.UTC)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeShift,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.ShiftRecurrenceRule{
					NumberOfTaskDays:  2,
					NumberOfShiftDays: 2,
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(106),
			Date:        domaintask.Date(time.Date(2026, time.June, 1, 14, 0, 0, 0, time.UTC)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeParity,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.ParityRecurrenceRule{
					IsEven: true,
				},
			},
		},
	}

	generated := generateTasksForPeriod(tasks, from, to)
	gotCountByID := make(map[int]int)
	for _, generatedTask := range generated {
		gotCountByID[generatedTask.ID.Int()]++
	}

	wantCountByID := map[int]int{
		101: 5,  // weekly Mondays: 1,8,15,22,29
		102: 1,  // monthly day 15
		103: 1,  // yearly 06-07
		104: 3,  // odd ISO weeks Mondays: 1,15,29
		105: 16, // shift 2 on / 2 off in June
		106: 15, // even dates in June
	}

	for taskID, expectedCount := range wantCountByID {
		if gotCountByID[taskID] != expectedCount {
			t.Fatalf("task id %d: expected %d occurrences, got %d", taskID, expectedCount, gotCountByID[taskID])
		}
	}
}

func TestTaskServiceGetTasks_FiltersDoneAndExpired(t *testing.T) {
	now := time.Now().UTC()
	from := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -2)
	to := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.UTC).AddDate(0, 0, 2)

	expiredDate := from.Add(12 * time.Hour)
	doneDate := from.AddDate(0, 0, 1).Add(12 * time.Hour)
	plannedDate := to.AddDate(0, 0, -1).Add(-12 * time.Hour)

	repo := &taskRepoMock{
		getTasksForPeriodFn: func(ctx context.Context, qFrom time.Time, qTo time.Time) ([]*domaintask.Task, error) {
			return []*domaintask.Task{
				{ID: domaintask.IDFromInt(1), Date: domaintask.Date(expiredDate), IsRecurring: false, IsDone: false},
				{ID: domaintask.IDFromInt(2), Date: domaintask.Date(doneDate), IsRecurring: false, IsDone: true},
				{ID: domaintask.IDFromInt(3), Date: domaintask.Date(plannedDate), IsRecurring: false, IsDone: false},
			}, nil
		},
	}
	svc := &TaskService{repo: repo}

	got, err := svc.GetTasks(context.Background(), from, to, []domaintask.TaskStatus{
		domaintask.TaskStatusDone,
		domaintask.TaskStatusExpired,
	})
	if err != nil {
		t.Fatalf("GetTasks returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks after status filtering, got %d", len(got))
	}

	statuses := map[domaintask.TaskStatus]bool{}
	for _, task := range got {
		statuses[task.Status] = true
	}
	if !statuses[domaintask.TaskStatusDone] || !statuses[domaintask.TaskStatusExpired] {
		t.Fatalf("expected statuses to contain done and expired, got %+v", statuses)
	}
	if statuses[domaintask.TaskStatusPlanned] {
		t.Fatalf("expected planned status to be filtered out")
	}
}

func TestTaskServiceMarkTaskDone_RecurringMarksOnlyOccurrence(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := time.Date(2026, time.June, 26, 9, 0, 0, 0, time.UTC)
	end := time.Date(2026, time.July, 10, 23, 59, 59, 0, time.UTC)
	occurrenceInput := time.Date(2026, time.June, 27, 18, 30, 0, 0, time.UTC)
	expectedOccurrenceDay := time.Date(2026, time.June, 27, 0, 0, 0, 0, time.UTC)

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

func TestTaskServiceGetTasks_RejectsZeroPeriodBounds(t *testing.T) {
	svc := &TaskService{repo: &taskRepoMock{}}

	_, err := svc.GetTasks(context.Background(), time.Time{}, time.Now().UTC(), nil)
	if err == nil || !strings.Contains(err.Error(), "from and to are required") {
		t.Fatalf("expected zero bounds validation error, got: %v", err)
	}

	_, err = svc.GetTasks(context.Background(), time.Now().UTC(), time.Time{}, nil)
	if err == nil || !strings.Contains(err.Error(), "from and to are required") {
		t.Fatalf("expected zero bounds validation error, got: %v", err)
	}
}

func TestTaskServiceGetTasks_RejectsFromAfterTo(t *testing.T) {
	svc := &TaskService{repo: &taskRepoMock{}}
	from := time.Date(2026, time.June, 27, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.June, 26, 23, 59, 59, 0, time.UTC)

	_, err := svc.GetTasks(context.Background(), from, to, nil)
	if err == nil || !strings.Contains(err.Error(), "from must be less than or equal to to") {
		t.Fatalf("expected from<=to validation error, got: %v", err)
	}
}

func TestTaskServiceMarkTaskDone_RecurringRequiresOccurrenceDate(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	repo := &taskRepoMock{
		getTaskFn: func(ctx context.Context, id domaintask.ID) (*domaintask.Task, error) {
			return &domaintask.Task{
				ID:          taskID,
				Date:        domaintask.Date(time.Date(2026, time.June, 26, 9, 0, 0, 0, time.UTC)),
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
	if err == nil || !strings.Contains(err.Error(), "occurrence_date is required for recurring task") {
		t.Fatalf("expected occurrence_date required error, got: %v", err)
	}
}

func TestTaskServiceMarkTaskDone_RecurringRejectsOccurrenceBeforeStart(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := time.Date(2026, time.June, 26, 9, 0, 0, 0, time.UTC)
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
	if err == nil || !strings.Contains(err.Error(), "earlier than task start date") {
		t.Fatalf("expected start-date validation error, got: %v", err)
	}
}

func TestTaskServiceMarkTaskDone_RecurringRejectsOccurrenceAfterEnd(t *testing.T) {
	taskID := domaintask.IDFromInt(22)
	start := time.Date(2026, time.June, 26, 9, 0, 0, 0, time.UTC)
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
	if err == nil || !strings.Contains(err.Error(), "after recurrence end date") {
		t.Fatalf("expected end-date validation error, got: %v", err)
	}
}

func TestTaskServiceMarkTaskDone_RecurringWithoutEndDateAllowsFutureOccurrence(t *testing.T) {
	taskID := domaintask.IDFromInt(220)
	start := time.Date(2026, time.June, 22, 9, 0, 0, 0, time.UTC) // Monday
	occurrence := time.Date(2028, time.June, 19, 12, 0, 0, 0, time.UTC)

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
						WeekDay: 1, // Monday
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
	start := time.Date(2026, time.June, 22, 9, 0, 0, 0, time.UTC) // Monday
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
						WeekDay: 1, // Monday
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
	occurrence := time.Date(2026, time.June, 23, 13, 0, 0, 0, time.UTC) // Tuesday

	err := svc.MarkTaskDone(context.Background(), taskID, &occurrence)
	if err == nil || !strings.Contains(err.Error(), "does not match recurrence rule") {
		t.Fatalf("expected recurrence-rule mismatch error, got: %v", err)
	}
}

package task

import (
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func TestGenerateTasksForPeriod_IncludesFromAndToBoundaries(t *testing.T) {
	from := mustUTC(2026, time.June, 26, 0, 0, 0)
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
	from := mustUTC(2026, time.June, 26, 0, 0, 0)
	to := time.Date(2026, time.July, 3, 23, 59, 59, 999999999, time.UTC)
	start := mustUTC(2026, time.June, 26, 14, 4, 19)
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
	from := mustUTC(2026, time.June, 26, 0, 0, 0)
	to := time.Date(2026, time.June, 30, 23, 59, 59, 999999999, time.UTC)
	start := mustUTC(2026, time.June, 26, 9, 0, 0)

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
	if !time.Time(got[len(got)-1].Date).Equal(mustUTC(2026, time.June, 30, 9, 0, 0)) {
		t.Fatalf("expected last generated date to match query upper bound day")
	}
}

func TestGenerateTasksForPeriod_AllRecurringTypes(t *testing.T) {
	from := mustUTC(2026, time.June, 1, 0, 0, 0)
	to := time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)

	tasks := []*domaintask.Task{
		{
			ID:          domaintask.IDFromInt(101),
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 9, 0, 0)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeWeekly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.WeeklyRecurrenceRule{
					WeekDay: 1,
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(102),
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 10, 0, 0)),
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
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 11, 0, 0)),
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
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 12, 0, 0)),
			IsRecurring: true,
			Recurring: &domaintask.RecurringTask{
				RecurringType: domaintask.RecurringTypeBiweekly,
				EndDate:       recurringEndDate(time.Date(2026, time.June, 30, 23, 59, 59, 0, time.UTC)),
				Rule: domaintask.BiweeklyRecurrenceRule{
					IsOdd:   true,
					WeekDay: 1,
				},
			},
		},
		{
			ID:          domaintask.IDFromInt(105),
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 13, 0, 0)),
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
			Date:        domaintask.Date(mustUTC(2026, time.June, 1, 14, 0, 0)),
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
		101: 5,
		102: 1,
		103: 1,
		104: 3,
		105: 16,
		106: 15,
	}

	for taskID, expectedCount := range wantCountByID {
		if gotCountByID[taskID] != expectedCount {
			t.Fatalf("task id %d: expected %d occurrences, got %d", taskID, expectedCount, gotCountByID[taskID])
		}
	}
}

func TestGenerateTasksForPeriod_RecurringSetsIsDoneFromDoneDates(t *testing.T) {
	from := mustUTC(2026, time.June, 26, 0, 0, 0)
	to := time.Date(2026, time.June, 28, 23, 59, 59, 999999999, time.UTC)
	start := mustUTC(2026, time.June, 26, 9, 0, 0)

	recurring := &domaintask.Task{
		ID:          domaintask.IDFromInt(333),
		Date:        domaintask.Date(start),
		IsRecurring: true,
		Recurring: &domaintask.RecurringTask{
			RecurringType: domaintask.RecurringTypeShift,
			Rule: domaintask.ShiftRecurrenceRule{
				NumberOfTaskDays:  1,
				NumberOfShiftDays: 0,
			},
		},
		DoneDates: map[string]struct{}{
			"2026-06-27": {},
		},
	}

	got := generateTasksForPeriod([]*domaintask.Task{recurring}, from, to)
	if len(got) != 3 {
		t.Fatalf("expected 3 generated tasks, got %d", len(got))
	}

	if got[0].IsDone {
		t.Fatalf("expected first occurrence to be not done")
	}
	if !got[1].IsDone {
		t.Fatalf("expected second occurrence to be done")
	}
	if got[2].IsDone {
		t.Fatalf("expected third occurrence to be not done")
	}
}

func TestGenerateTasksForPeriod_SortsByDateThenID(t *testing.T) {
	from := mustUTC(2026, time.June, 26, 0, 0, 0)
	to := time.Date(2026, time.June, 26, 23, 59, 59, 0, time.UTC)

	tasks := []*domaintask.Task{
		{ID: domaintask.IDFromInt(9), Date: domaintask.Date(mustUTC(2026, time.June, 26, 12, 0, 0))},
		{ID: domaintask.IDFromInt(2), Date: domaintask.Date(mustUTC(2026, time.June, 26, 12, 0, 0))},
		{ID: domaintask.IDFromInt(4), Date: domaintask.Date(mustUTC(2026, time.June, 26, 7, 0, 0))},
	}

	got := generateTasksForPeriod(tasks, from, to)
	if len(got) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(got))
	}

	gotOrder := []int{got[0].ID.Int(), got[1].ID.Int(), got[2].ID.Int()}
	wantOrder := []int{4, 2, 9}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("expected id order %v, got %v", wantOrder, gotOrder)
		}
	}
}

package task

import (
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func TestRecurringRangeForPeriod_ClampsByTaskStartAndEnd(t *testing.T) {
	task := &domaintask.Task{
		Date: domaintask.Date(mustUTC(2026, time.June, 10, 9, 0, 0)),
		Recurring: &domaintask.RecurringTask{
			EndDate: recurringEndDate(mustUTC(2026, time.June, 20, 23, 59, 59)),
			Rule:    domaintask.WeeklyRecurrenceRule{WeekDay: 3},
		},
	}

	start, end, ok := recurringRangeForPeriod(task, mustUTC(2026, time.June, 1, 0, 0, 0), mustUTC(2026, time.June, 30, 0, 0, 0))
	if !ok {
		t.Fatalf("expected recurrence range to be valid")
	}
	if !start.Equal(mustUTC(2026, time.June, 10, 9, 0, 0)) {
		t.Fatalf("expected start to clamp to task date, got %s", start.Format(time.RFC3339))
	}
	if !end.Equal(mustUTC(2026, time.June, 20, 23, 59, 59)) {
		t.Fatalf("expected end to clamp to recurrence end, got %s", end.Format(time.RFC3339))
	}
}

func TestRecurringRangeForPeriod_ReturnsNotOkWhenStartAfterEnd(t *testing.T) {
	task := &domaintask.Task{
		Date: domaintask.Date(mustUTC(2026, time.June, 25, 9, 0, 0)),
		Recurring: &domaintask.RecurringTask{
			EndDate: recurringEndDate(mustUTC(2026, time.June, 26, 23, 59, 59)),
			Rule:    domaintask.WeeklyRecurrenceRule{WeekDay: 4},
		},
	}

	_, _, ok := recurringRangeForPeriod(task, mustUTC(2026, time.June, 27, 0, 0, 0), mustUTC(2026, time.June, 30, 0, 0, 0))
	if ok {
		t.Fatalf("expected recurrence range to be invalid")
	}
}

func TestMatchesRecurringRule_ShiftWithInvalidCycleReturnsFalse(t *testing.T) {
	task := &domaintask.Task{
		Date:        domaintask.Date(mustUTC(2026, time.June, 1, 9, 0, 0)),
		IsRecurring: true,
		Recurring: &domaintask.RecurringTask{
			RecurringType: domaintask.RecurringTypeShift,
			Rule: domaintask.ShiftRecurrenceRule{
				NumberOfTaskDays:  0,
				NumberOfShiftDays: 0,
			},
		},
	}
	if matchesRecurringRule(mustUTC(2026, time.June, 1, 9, 0, 0), task) {
		t.Fatalf("expected shift recurrence with invalid cycle to not match")
	}
}

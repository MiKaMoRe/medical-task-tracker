package task

import (
	"testing"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func TestValidateTask_RecurringPayloadMustMatchFlag(t *testing.T) {
	t.Parallel()

	rule := domaintask.WeeklyRecurrenceRule{WeekDay: 1}
	recurring := &domaintask.RecurringTask{
		RecurringType: domaintask.RecurringTypeWeekly,
		Rule:          rule,
	}

	tests := []struct {
		name       string
		task       *domaintask.Task
		wantErrSub string
	}{
		{
			name: "requires recurring payload when recurring is enabled",
			task: &domaintask.Task{
				IsRecurring: true,
				Recurring:   nil,
			},
			wantErrSub: "recurring task is required when is_recurring is true",
		},
		{
			name: "rejects recurring payload when recurring is disabled",
			task: &domaintask.Task{
				IsRecurring: false,
				Recurring:   recurring,
			},
			wantErrSub: "recurring payload must be empty when is_recurring is false",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateTask(tc.task).Err()
			assertErrorContains(t, err, tc.wantErrSub)
		})
	}
}

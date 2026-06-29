package task

import (
	"testing"
	"time"

	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/repository"
)

func TestToCreateTaskDTO_DropsRecurringPayloadWhenFlagIsFalse(t *testing.T) {
	t.Parallel()

	rule := domaintask.WeeklyRecurrenceRule{WeekDay: 1}
	domainTask := &domaintask.Task{
		Title:       domaintask.Title("visit doctor"),
		Description: domaintask.Description("follow-up"),
		Date:        domaintask.Date(time.Date(2026, time.June, 29, 10, 0, 0, 0, time.UTC)),
		IsRecurring: false,
		Recurring: &domaintask.RecurringTask{
			RecurringType: domaintask.RecurringTypeWeekly,
			Rule:          rule,
		},
	}

	dto := toCreateTaskDTO(domainTask)
	if dto.RecurringRule != nil {
		t.Fatalf("expected recurring payload to be empty when is_recurring is false")
	}
}

func TestToDomainTask_DropsRecurringPayloadWhenFlagIsFalse(t *testing.T) {
	t.Parallel()

	dto := &TaskDTO{
		Base:         repository.Base{ID: 11},
		Name:         "visit doctor",
		Description:  "follow-up",
		Date:         time.Date(2026, time.June, 29, 10, 0, 0, 0, time.UTC),
		IsRecurrence: false,
		RecurrenceRule: &RecurrenceRuleDTO{
			RuleType: string(domaintask.RecurringTypeWeekly),
			WeeklyRule: &RecurrenceWeeklyRuleDTO{
				WeekDay: 1,
			},
		},
	}

	task := toDomainTask(dto)
	if task.Recurring != nil {
		t.Fatalf("expected recurring payload to be empty when is_recurring is false")
	}
}

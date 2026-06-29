package task

import (
	"strconv"
	"time"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func toCreateTaskDTO(t *task.Task) CreateTaskDTO {
	dto := CreateTaskDTO{
		Title:       string(t.Title),
		Description: string(t.Description),
		Date:        time.Time(t.Date),
		IsRecurring: t.IsRecurring,
		Tags:        tagNames(t.Tags),
	}

	if t.IsRecurring && t.Recurring != nil {
		recurringRuleDTO := &CreateRecurringRuleDTO{
			RecurringType: t.Recurring.RecurringType.String(),
		}
		if t.Recurring.EndDate != nil {
			endDate := time.Time(*t.Recurring.EndDate)
			recurringRuleDTO.EndDate = &endDate
		}
		switch rule := t.Recurring.Rule.(type) {
		case task.WeeklyRecurrenceRule:
			recurringRuleDTO.WeeklyRule = &CreateRecurrenceWeeklyRuleDTO{
				WeekDay: rule.WeekDay,
			}
		case task.MonthlyRecurrenceRule:
			recurringRuleDTO.MonthlyRule = &CreateRecurrenceMonthlyRuleDTO{
				MonthDay: rule.MonthDay,
			}
		case task.YearlyRecurrenceRule:
			recurringRuleDTO.YearlyRule = &CreateRecurrenceYearlyRuleDTO{
				Month: rule.Month,
				Day:   rule.Day,
			}
		case task.BiweeklyRecurrenceRule:
			recurringRuleDTO.BiweeklyRule = &CreateRecurrenceBiweeklyRuleDTO{
				IsOdd:   rule.IsOdd,
				WeekDay: rule.WeekDay,
			}
		case task.ShiftRecurrenceRule:
			recurringRuleDTO.ShiftRule = &CreateRecurrenceShiftRuleDTO{
				NumberOfTaskDays:  rule.NumberOfTaskDays,
				NumberOfShiftDays: rule.NumberOfShiftDays,
			}
		case task.ParityRecurrenceRule:
			recurringRuleDTO.ParityRule = &CreateRecurrenceParityRuleDTO{
				IsEven: rule.IsEven,
			}
		}
		dto.RecurringRule = recurringRuleDTO
	}

	return dto
}

func toTaskDTO(dto CreateTaskDTO) *TaskDTO {
	taskDTO := &TaskDTO{
		Name:         dto.Title,
		Description:  dto.Description,
		Date:         dto.Date,
		IsRecurrence: dto.IsRecurring,
	}

	if dto.RecurringRule != nil {
		recurrenceRuleDTO := &RecurrenceRuleDTO{
			RuleType: dto.RecurringRule.RecurringType,
			EndDate:  dto.RecurringRule.EndDate,
		}
		if dto.RecurringRule.WeeklyRule != nil {
			recurrenceRuleDTO.WeeklyRule = &RecurrenceWeeklyRuleDTO{
				WeekDay: dto.RecurringRule.WeeklyRule.WeekDay,
			}
		}
		if dto.RecurringRule.MonthlyRule != nil {
			recurrenceRuleDTO.MonthlyRule = &RecurrenceMonthlyRuleDTO{
				MonthDay: dto.RecurringRule.MonthlyRule.MonthDay,
			}
		}
		if dto.RecurringRule.YearlyRule != nil {
			recurrenceRuleDTO.YearlyRule = &RecurrenceYearlyRuleDTO{
				Month: dto.RecurringRule.YearlyRule.Month,
				Day:   dto.RecurringRule.YearlyRule.Day,
			}
		}
		if dto.RecurringRule.BiweeklyRule != nil {
			recurrenceRuleDTO.BiweeklyRule = &RecurrenceBiweeklyRuleDTO{
				IsOdd:   dto.RecurringRule.BiweeklyRule.IsOdd,
				WeekDay: dto.RecurringRule.BiweeklyRule.WeekDay,
			}
		}
		if dto.RecurringRule.ShiftRule != nil {
			recurrenceRuleDTO.ShiftRule = &RecurrenceShiftRuleDTO{
				NumberOfTaskDays:  dto.RecurringRule.ShiftRule.NumberOfTaskDays,
				NumberOfShiftDays: dto.RecurringRule.ShiftRule.NumberOfShiftDays,
			}
		}
		if dto.RecurringRule.ParityRule != nil {
			recurrenceRuleDTO.ParityRule = &RecurrenceParityRuleDTO{
				IsEven: dto.RecurringRule.ParityRule.IsEven,
			}
		}
		taskDTO.RecurrenceRule = recurrenceRuleDTO
	}

	if len(dto.Tags) > 0 {
		taskDTO.Tags = make([]TagDTO, len(dto.Tags))
		for i, name := range dto.Tags {
			taskDTO.Tags[i] = TagDTO{Name: name}
		}
	}

	return taskDTO
}

func toDomainTask(dto *TaskDTO) *task.Task {
	t := &task.Task{
		ID:          task.IDFromInt(dto.ID),
		Title:       task.Title(dto.Name),
		Description: task.Description(dto.Description),
		Date:        task.Date(dto.Date),
		IsRecurring: dto.IsRecurrence,
		Tags:        toDomainTags(dto.Tags),
		IsDone:      dto.IsDone,
		DoneDates:   toDoneDatesMap(dto.Completions),
	}

	if dto.IsRecurrence && dto.RecurrenceRule != nil {
		recurringTask := &task.RecurringTask{
			RecurringType: task.RecurringType(dto.RecurrenceRule.RuleType),
			Rule:          mapRecurrenceRule(dto.RecurrenceRule),
		}
		if dto.RecurrenceRule.EndDate != nil {
			endDate := task.Date(*dto.RecurrenceRule.EndDate)
			recurringTask.EndDate = &endDate
		}
		t.Recurring = recurringTask
	}

	return t
}

func toDoneDatesMap(dtos []TaskOccurrenceCompletionDTO) map[string]struct{} {
	doneDates := make(map[string]struct{}, len(dtos))
	for _, dto := range dtos {
		doneDates[dto.OccurrenceDate.UTC().Format(time.DateOnly)] = struct{}{}
	}
	return doneDates
}

func mapRecurrenceRule(dto *RecurrenceRuleDTO) task.RecurrenceRule {
	if dto == nil {
		return nil
	}

	switch dto.RuleType {
	case string(task.RecurringTypeWeekly):
		if dto.WeeklyRule == nil {
			return nil
		}
		return task.WeeklyRecurrenceRule{
			WeekDay: dto.WeeklyRule.WeekDay,
		}
	case string(task.RecurringTypeMonthly):
		if dto.MonthlyRule == nil {
			return nil
		}
		return task.MonthlyRecurrenceRule{
			MonthDay: dto.MonthlyRule.MonthDay,
		}
	case string(task.RecurringTypeYearly):
		if dto.YearlyRule == nil {
			return nil
		}
		return task.YearlyRecurrenceRule{
			Month: dto.YearlyRule.Month,
			Day:   dto.YearlyRule.Day,
		}
	case string(task.RecurringTypeBiweekly):
		if dto.BiweeklyRule == nil {
			return nil
		}
		return task.BiweeklyRecurrenceRule{
			IsOdd:   dto.BiweeklyRule.IsOdd,
			WeekDay: dto.BiweeklyRule.WeekDay,
		}
	case string(task.RecurringTypeShift):
		if dto.ShiftRule == nil {
			return nil
		}
		return task.ShiftRecurrenceRule{
			NumberOfTaskDays:  dto.ShiftRule.NumberOfTaskDays,
			NumberOfShiftDays: dto.ShiftRule.NumberOfShiftDays,
		}
	case string(task.RecurringTypeParity):
		if dto.ParityRule == nil {
			return nil
		}
		return task.ParityRecurrenceRule{
			IsEven: dto.ParityRule.IsEven,
		}
	default:
		return nil
	}
}

func toDomainTags(dtos []TagDTO) []tag.Tag {
	tags := make([]tag.Tag, len(dtos))
	for i, dto := range dtos {
		tags[i] = tag.Tag{
			ID:   tag.ID(strconv.Itoa(dto.ID)),
			Name: tag.Name(dto.Name),
		}
	}
	return tags
}

func tagNames(tags []tag.Tag) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name.String()
	}
	return names
}

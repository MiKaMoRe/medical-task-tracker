package task

import (
	"errors"
	"fmt"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
)

type RecurringType string

const (
	RecurringTypeWeekly   RecurringType = "weekly"
	RecurringTypeMonthly  RecurringType = "monthly"
	RecurringTypeYearly   RecurringType = "yearly"
	RecurringTypeBiweekly RecurringType = "biweekly"
	RecurringTypeShift    RecurringType = "shift"
	RecurringTypeParity   RecurringType = "parity"
)

func (r RecurringType) String() string {
	return string(r)
}

func NewRecurringType(s string) (RecurringType, error) {
	switch s {
	case "weekly":
		return RecurringTypeWeekly, nil
	case "monthly":
		return RecurringTypeMonthly, nil
	case "yearly":
		return RecurringTypeYearly, nil
	case "biweekly":
		return RecurringTypeBiweekly, nil
	case "shift":
		return RecurringTypeShift, nil
	case "parity":
		return RecurringTypeParity, nil
	default:
		return RecurringType(s), fmt.Errorf("invalid recurring type: %s", s)
	}
}

type RecurringTask struct {
	RecurringType RecurringType  `json:"recurring_type"`
	EndDate       *Date          `json:"end_date,omitempty"`
	Rule          RecurrenceRule `json:"-"`
}

type RecurrenceRule interface {
	isRecurrenceRule()
}

type RecurrenceRuleInput struct {
	Weekly   *WeeklyRecurrenceRuleInput
	Monthly  *MonthlyRecurrenceRuleInput
	Yearly   *YearlyRecurrenceRuleInput
	Biweekly *BiweeklyRecurrenceRuleInput
	Shift    *ShiftRecurrenceRuleInput
	Parity   *ParityRecurrenceRuleInput
}

type WeeklyRecurrenceRuleInput struct {
	WeekDay *int
}

type MonthlyRecurrenceRuleInput struct {
	MonthDay *int
}

type YearlyRecurrenceRuleInput struct {
	Month *int
	Day   *int
}

type BiweeklyRecurrenceRuleInput struct {
	IsOdd   *bool
	WeekDay *int
}

type ShiftRecurrenceRuleInput struct {
	NumberOfTaskDays  *int
	NumberOfShiftDays *int
}

type ParityRecurrenceRuleInput struct {
	IsEven *bool
}

type WeeklyRecurrenceRule struct {
	WeekDay int
}

func (WeeklyRecurrenceRule) isRecurrenceRule() {}

type MonthlyRecurrenceRule struct {
	MonthDay int
}

func (MonthlyRecurrenceRule) isRecurrenceRule() {}

type YearlyRecurrenceRule struct {
	Month int
	Day   int
}

func (YearlyRecurrenceRule) isRecurrenceRule() {}

type BiweeklyRecurrenceRule struct {
	IsOdd   bool
	WeekDay int
}

func (BiweeklyRecurrenceRule) isRecurrenceRule() {}

type ShiftRecurrenceRule struct {
	NumberOfTaskDays  int
	NumberOfShiftDays int
}

func (ShiftRecurrenceRule) isRecurrenceRule() {}

type ParityRecurrenceRule struct {
	IsEven bool
}

func (ParityRecurrenceRule) isRecurrenceRule() {}

type recurrenceRuleFactory func(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap)

var recurrenceRuleFactories = map[RecurringType]recurrenceRuleFactory{
	RecurringTypeWeekly:   buildWeeklyRule,
	RecurringTypeMonthly:  buildMonthlyRule,
	RecurringTypeYearly:   buildYearlyRule,
	RecurringTypeBiweekly: buildBiweeklyRule,
	RecurringTypeShift:    buildShiftRule,
	RecurringTypeParity:   buildParityRule,
}

func NewRecurringTask(recurringType string, endDateT *time.Time, ruleInput RecurrenceRuleInput) (*RecurringTask, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	var endDate *Date
	if endDateT != nil {
		if endDateT.Before(time.Now().In(time.UTC)) {
			vm.Add("recurring.end_date", errors.New("end date must be after current date"))
		}

		normalizedEndDate, errs := NewDate(*endDateT)
		vm.Add("recurring.end_date", errs...)
		endDate = &normalizedEndDate
	}

	rType, err := NewRecurringType(recurringType)
	vm.Add("recurring.type", err)

	var recurrenceRule RecurrenceRule
	if err == nil {
		var ruleErrs *apperrors.ValidationMap
		recurrenceRule, ruleErrs = buildRecurrenceRule(rType, ruleInput)
		vm = apperrors.MergeValidationMaps(vm, ruleErrs)
	}

	return &RecurringTask{
		RecurringType: rType,
		EndDate:       endDate,
		Rule:          recurrenceRule,
	}, vm
}

func buildRecurrenceRule(recurringType RecurringType, input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	factory, ok := recurrenceRuleFactories[recurringType]
	if !ok {
		vm := apperrors.NewValidationMap()
		vm.Add("recurring.type", fmt.Errorf("unsupported recurring type: %s", recurringType))
		return nil, vm
	}
	return factory(input)
}

func buildWeeklyRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Weekly == nil {
		vm.Add("recurring.weekly_rule", errors.New("weekly rule is required"))
		return nil, vm
	}
	weekDay, ok := requiredInt("recurring.weekly_rule.week_day", input.Weekly.WeekDay, vm)
	if ok && (weekDay < 0 || weekDay > 6) {
		vm.Add("recurring.weekly_rule.week_day", errors.New("week day must be between 0 and 6"))
	}
	return WeeklyRecurrenceRule{WeekDay: weekDay}, vm
}

func buildMonthlyRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Monthly == nil {
		vm.Add("recurring.monthly_rule", errors.New("monthly rule is required"))
		return nil, vm
	}
	monthDay, ok := requiredInt("recurring.monthly_rule.month_day", input.Monthly.MonthDay, vm)
	if ok && (monthDay < 1 || monthDay > 31) {
		vm.Add("recurring.monthly_rule.month_day", errors.New("month day must be between 1 and 31"))
	}
	return MonthlyRecurrenceRule{MonthDay: monthDay}, vm
}

func buildYearlyRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Yearly == nil {
		vm.Add("recurring.yearly_rule", errors.New("yearly rule is required"))
		return nil, vm
	}
	month, monthOK := requiredInt("recurring.yearly_rule.month", input.Yearly.Month, vm)
	if monthOK && (month < 1 || month > 12) {
		vm.Add("recurring.yearly_rule.month", errors.New("month must be between 1 and 12"))
	}
	day, dayOK := requiredInt("recurring.yearly_rule.day", input.Yearly.Day, vm)
	if dayOK && (day < 1 || day > 31) {
		vm.Add("recurring.yearly_rule.day", errors.New("day must be between 1 and 31"))
	}
	return YearlyRecurrenceRule{Month: month, Day: day}, vm
}

func buildBiweeklyRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Biweekly == nil {
		vm.Add("recurring.biweekly_rule", errors.New("biweekly rule is required"))
		return nil, vm
	}
	isOdd, _ := requiredBool("recurring.biweekly_rule.is_odd", input.Biweekly.IsOdd, vm)
	weekDay, weekDayOK := requiredInt("recurring.biweekly_rule.week_day", input.Biweekly.WeekDay, vm)
	if weekDayOK && (weekDay < 0 || weekDay > 6) {
		vm.Add("recurring.biweekly_rule.week_day", errors.New("week day must be between 0 and 6"))
	}
	return BiweeklyRecurrenceRule{IsOdd: isOdd, WeekDay: weekDay}, vm
}

func buildShiftRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Shift == nil {
		vm.Add("recurring.shift_rule", errors.New("shift rule is required"))
		return nil, vm
	}
	numberOfTaskDays, taskDaysOK := requiredInt("recurring.shift_rule.number_of_task_days", input.Shift.NumberOfTaskDays, vm)
	if taskDaysOK && numberOfTaskDays <= 0 {
		vm.Add("recurring.shift_rule.number_of_task_days", errors.New("number of task days must be greater than 0"))
	}
	numberOfShiftDays, shiftDaysOK := requiredInt("recurring.shift_rule.number_of_shift_days", input.Shift.NumberOfShiftDays, vm)
	if shiftDaysOK && numberOfShiftDays < 0 {
		vm.Add("recurring.shift_rule.number_of_shift_days", errors.New("number of shift days must be greater than or equal to 0"))
	}
	return ShiftRecurrenceRule{
		NumberOfTaskDays:  numberOfTaskDays,
		NumberOfShiftDays: numberOfShiftDays,
	}, vm
}

func buildParityRule(input RecurrenceRuleInput) (RecurrenceRule, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	if input.Parity == nil {
		vm.Add("recurring.parity_rule", errors.New("parity rule is required"))
		return nil, vm
	}
	isEven, _ := requiredBool("recurring.parity_rule.is_even", input.Parity.IsEven, vm)
	return ParityRecurrenceRule{IsEven: isEven}, vm
}

func requiredInt(field string, value *int, vm *apperrors.ValidationMap) (int, bool) {
	if value == nil {
		vm.Add(field, errors.New("field is required"))
		return 0, false
	}
	return *value, true
}

func requiredBool(field string, value *bool, vm *apperrors.ValidationMap) (bool, bool) {
	if value == nil {
		vm.Add(field, errors.New("field is required"))
		return false, false
	}
	return *value, true
}

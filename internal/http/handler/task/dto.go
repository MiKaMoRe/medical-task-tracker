package task

import "time"

type CreateTaskRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        time.Time            `json:"date"`
	IsRecurring bool                 `json:"is_recurring"`
	Tags        []string             `json:"tags"`
	Recurring   *CreateRecurringTask `json:"recurring,omitempty"`
}

type UpdateTaskRequest struct {
	Title       string               `json:"title"`
	Description string               `json:"description"`
	Date        time.Time            `json:"date"`
	IsRecurring bool                 `json:"is_recurring"`
	Tags        []string             `json:"tags"`
	Recurring   *CreateRecurringTask `json:"recurring,omitempty"`
}

type CreateRecurringTask struct {
	RecurringType string                     `json:"recurring_type"`
	EndDate       *time.Time                 `json:"end_date,omitempty"`
	WeeklyRule    *CreateWeeklyRuleRequest   `json:"recurrence_weekly_rule,omitempty"`
	MonthlyRule   *CreateMonthlyRuleRequest  `json:"recurrence_monthly_rule,omitempty"`
	YearlyRule    *CreateYearlyRuleRequest   `json:"recurrence_yearly_rule,omitempty"`
	BiweeklyRule  *CreateBiweeklyRuleRequest `json:"recurrence_biweekly_rule,omitempty"`
	ShiftRule     *CreateShiftRuleRequest    `json:"recurrence_shift_rule,omitempty"`
	ParityRule    *CreateParityRuleRequest   `json:"recurrence_parity_rule,omitempty"`
}

type CreateWeeklyRuleRequest struct {
	WeekDay *int `json:"week_day"`
}

type CreateMonthlyRuleRequest struct {
	MonthDay *int `json:"month_day"`
}

type CreateYearlyRuleRequest struct {
	Month *int `json:"month"`
	Day   *int `json:"day"`
}

type CreateBiweeklyRuleRequest struct {
	IsOdd   *bool `json:"is_odd"`
	WeekDay *int  `json:"week_day"`
}

type CreateShiftRuleRequest struct {
	NumberOfTaskDays  *int `json:"number_of_task_days"`
	NumberOfShiftDays *int `json:"number_of_shift_days"`
}

type CreateParityRuleRequest struct {
	IsEven *bool `json:"is_even"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	DueDate     string    `json:"due_date"`
	IsRecurring bool      `json:"is_recurring"`
	Tags        []string  `json:"tags"`
}

type MarkTaskDoneRequest struct {
	OccurrenceDate *string `json:"occurrence_date,omitempty"`
}

type TaskTagsRequest struct {
	Tags []string `json:"tags"`
}

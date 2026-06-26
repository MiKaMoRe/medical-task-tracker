package task

import (
	"time"

	domaintag "github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

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
	ID int `json:"id"`
}

type GetTaskRequest struct {
	ID string `json:"id"`
}

type GetTaskResponse struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Date        time.Time              `json:"date"`
	IsRecurring bool                   `json:"is_recurring"`
	Status      string                 `json:"status"`
	Recurring   *RecurringTaskResponse `json:"recurring,omitempty"`
	Tags        []string               `json:"tags"`
}

type MarkTaskDoneRequest struct {
	OccurrenceDate *string `json:"occurrence_date,omitempty"`
}

type TaskTagsRequest struct {
	Tags []string `json:"tags"`
}

type TaskResponse struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Date        time.Time              `json:"date"`
	IsRecurring bool                   `json:"is_recurring"`
	Status      string                 `json:"status"`
	Recurring   *RecurringTaskResponse `json:"recurring,omitempty"`
	Tags        []TagResponse          `json:"tags"`
}

type RecurringTaskResponse struct {
	RecurringType string     `json:"recurring_type"`
	EndDate       *time.Time `json:"end_date,omitempty"`
}

type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func mapTaskResponse(t *domaintask.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID.Int(),
		Title:       string(t.Title),
		Description: string(t.Description),
		Date:        time.Time(t.Date),
		IsRecurring: t.IsRecurring,
		Status:      string(t.Status),
		Recurring:   mapRecurringResponse(t.Recurring),
		Tags:        mapTagResponses(t.Tags),
	}
}

func mapTaskResponses(tasks []*domaintask.Task) []TaskResponse {
	items := make([]TaskResponse, len(tasks))
	for i := range tasks {
		items[i] = mapTaskResponse(tasks[i])
	}
	return items
}

func mapRecurringResponse(recurring *domaintask.RecurringTask) *RecurringTaskResponse {
	if recurring == nil {
		return nil
	}

	var endDate *time.Time
	if recurring.EndDate != nil {
		value := time.Time(*recurring.EndDate)
		endDate = &value
	}

	return &RecurringTaskResponse{
		RecurringType: recurring.RecurringType.String(),
		EndDate:       endDate,
	}
}

func mapTagResponses(tags []domaintag.Tag) []TagResponse {
	items := make([]TagResponse, len(tags))
	for i := range tags {
		items[i] = TagResponse{
			ID:   tags[i].ID.String(),
			Name: tags[i].Name.String(),
		}
	}
	return items
}

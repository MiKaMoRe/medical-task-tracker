package task

import (
	"encoding/json"
	"net/http"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	taskDomain "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodPost) {
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid request body"))
		return
	}

	task, err := mapTaskRequest(req.Title, req.Description, req.Date, req.IsRecurring, req.Tags, req.Recurring)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	createdTask, err := h.srvc.CreateTask(r.Context(), task)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Created(w, createdTask)
}

func mapTaskRequest(
	title string,
	description string,
	date time.Time,
	isRecurring bool,
	tags []string,
	recurringReq *CreateRecurringTask,
) (*task.Task, error) {
	valErrs := apperrors.NewValidationMap()

	task, errs := taskDomain.NewTask(
		title,
		description,
		date,
		isRecurring,
		tags,
	)
	if !errs.IsEmpty() {
		valErrs = apperrors.MergeValidationMaps(valErrs, errs)
	}

	if recurringReq != nil {
		recurring, errs := taskDomain.NewRecurringTask(
			recurringReq.RecurringType,
			recurringReq.EndDate,
			mapRecurrenceRuleInput(recurringReq),
		)
		if !errs.IsEmpty() {
			valErrs = apperrors.MergeValidationMaps(valErrs, errs)
		}
		task.Recurring = recurring
	}
	return task, valErrs.Err()
}

func mapRecurrenceRuleInput(req *CreateRecurringTask) taskDomain.RecurrenceRuleInput {
	if req == nil {
		return taskDomain.RecurrenceRuleInput{}
	}

	return taskDomain.RecurrenceRuleInput{
		Weekly:   toWeeklyRuleInput(req.WeeklyRule),
		Monthly:  toMonthlyRuleInput(req.MonthlyRule),
		Yearly:   toYearlyRuleInput(req.YearlyRule),
		Biweekly: toBiweeklyRuleInput(req.BiweeklyRule),
		Shift:    toShiftRuleInput(req.ShiftRule),
		Parity:   toParityRuleInput(req.ParityRule),
	}
}

func toWeeklyRuleInput(req *CreateWeeklyRuleRequest) *taskDomain.WeeklyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.WeeklyRecurrenceRuleInput{
		WeekDay: req.WeekDay,
	}
}

func toMonthlyRuleInput(req *CreateMonthlyRuleRequest) *taskDomain.MonthlyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.MonthlyRecurrenceRuleInput{
		MonthDay: req.MonthDay,
	}
}

func toYearlyRuleInput(req *CreateYearlyRuleRequest) *taskDomain.YearlyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.YearlyRecurrenceRuleInput{
		Month: req.Month,
		Day:   req.Day,
	}
}

func toBiweeklyRuleInput(req *CreateBiweeklyRuleRequest) *taskDomain.BiweeklyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.BiweeklyRecurrenceRuleInput{
		IsOdd:   req.IsOdd,
		WeekDay: req.WeekDay,
	}
}

func toShiftRuleInput(req *CreateShiftRuleRequest) *taskDomain.ShiftRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.ShiftRecurrenceRuleInput{
		NumberOfTaskDays:  req.NumberOfTaskDays,
		NumberOfShiftDays: req.NumberOfShiftDays,
	}
}

func toParityRuleInput(req *CreateParityRuleRequest) *taskDomain.ParityRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &taskDomain.ParityRecurrenceRuleInput{
		IsEven: req.IsEven,
	}
}

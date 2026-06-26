package task

import (
	"encoding/json"
	"net/http"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
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

	if err := response.Created(w, mapTaskResponse(createdTask)); err != nil {
		h.handleError(w, r, err)
	}
}

func mapTaskRequest(
	title string,
	description string,
	date time.Time,
	isRecurring bool,
	tags []string,
	recurringReq *CreateRecurringTask,
) (*domaintask.Task, error) {
	valErrs := apperrors.NewValidationMap()

	task, errs := domaintask.NewTask(
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
		recurring, errs := domaintask.NewRecurringTask(
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

func mapRecurrenceRuleInput(req *CreateRecurringTask) domaintask.RecurrenceRuleInput {
	if req == nil {
		return domaintask.RecurrenceRuleInput{}
	}

	return domaintask.RecurrenceRuleInput{
		Weekly:   toWeeklyRuleInput(req.WeeklyRule),
		Monthly:  toMonthlyRuleInput(req.MonthlyRule),
		Yearly:   toYearlyRuleInput(req.YearlyRule),
		Biweekly: toBiweeklyRuleInput(req.BiweeklyRule),
		Shift:    toShiftRuleInput(req.ShiftRule),
		Parity:   toParityRuleInput(req.ParityRule),
	}
}

func toWeeklyRuleInput(req *CreateWeeklyRuleRequest) *domaintask.WeeklyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.WeeklyRecurrenceRuleInput{
		WeekDay: req.WeekDay,
	}
}

func toMonthlyRuleInput(req *CreateMonthlyRuleRequest) *domaintask.MonthlyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.MonthlyRecurrenceRuleInput{
		MonthDay: req.MonthDay,
	}
}

func toYearlyRuleInput(req *CreateYearlyRuleRequest) *domaintask.YearlyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.YearlyRecurrenceRuleInput{
		Month: req.Month,
		Day:   req.Day,
	}
}

func toBiweeklyRuleInput(req *CreateBiweeklyRuleRequest) *domaintask.BiweeklyRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.BiweeklyRecurrenceRuleInput{
		IsOdd:   req.IsOdd,
		WeekDay: req.WeekDay,
	}
}

func toShiftRuleInput(req *CreateShiftRuleRequest) *domaintask.ShiftRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.ShiftRecurrenceRuleInput{
		NumberOfTaskDays:  req.NumberOfTaskDays,
		NumberOfShiftDays: req.NumberOfShiftDays,
	}
}

func toParityRuleInput(req *CreateParityRuleRequest) *domaintask.ParityRecurrenceRuleInput {
	if req == nil {
		return nil
	}
	return &domaintask.ParityRecurrenceRuleInput{
		IsEven: req.IsEven,
	}
}

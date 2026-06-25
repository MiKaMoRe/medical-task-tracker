package task

import (
	"encoding/json"
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, err)
		return
	}

	task, err := mapTask(&req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	createdTask, err := h.srvc.CreateTask(r.Context(), task)
	if err != nil {
		h.handleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(createdTask)
}

func mapTask(req *CreateTaskRequest) (*task.Task, error) {
	title, err := task.NewTitle(req.Title)
	if err != nil {
		return nil, err
	}

	description, err := task.NewDescription(req.Description)
	if err != nil {
		return nil, err
	}

	date, err := task.NewDate(req.Date)
	if err != nil {
		return nil, err
	}

	tags, err := tag.NewTags(req.Tags)
	if err != nil {
		return nil, err
	}

	recurring, err := task.NewRecurringTask(req.Recurring.RecurringType, req.Date, req.Recurring.EndDate)
	if err != nil {
		return nil, err
	}

	return &task.Task{
		Title:       title,
		Description: description,
		Date:        date,
		IsRecurring: req.IsRecurring,
		Tags:        tags,
		Recurring:   recurring,
	}, nil
}

package task

import (
	"net/http"
	"strconv"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodGet) {
		return
	}

	statuses, err := parseStatusFilters(r.URL.Query()["status"])
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid task id"))
		return
	}
	task, err := h.srvc.GetTask(r.Context(), task.ID(id))
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	if len(statuses) > 0 && !containsTaskStatus(statuses, task.Status) {
		h.handleError(w, r, apperrors.NotFound("task not found"))
		return
	}
	response.Ok(w, task)
}

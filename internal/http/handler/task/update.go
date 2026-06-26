package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodPut) {
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid task id"))
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid request body"))
		return
	}

	updatedTask, err := mapTaskRequest(req.Title, req.Description, req.Date, req.IsRecurring, req.Tags, req.Recurring)
	if err != nil {
		h.handleError(w, r, err)
		return
	}
	updatedTask.ID = task.IDFromInt(id)

	savedTask, err := h.srvc.UpdateTask(r.Context(), updatedTask)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Ok(w, savedTask)
}

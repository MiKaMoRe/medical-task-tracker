package task

import (
	"net/http"
	"strconv"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
)

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodDelete) {
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid task id"))
		return
	}

	if err := h.srvc.DeleteTask(r.Context(), task.IDFromInt(id)); err != nil {
		h.handleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

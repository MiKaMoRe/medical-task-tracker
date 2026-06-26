package task

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) MarkTaskDone(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodPost) {
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid task id"))
		return
	}

	var req MarkTaskDoneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		h.handleError(w, r, apperrors.NewAppError("invalid request body"))
		return
	}

	var occurrenceDate *time.Time
	if req.OccurrenceDate != nil {
		normalized, err := parsePeriodDate(*req.OccurrenceDate, false)
		if err != nil {
			h.handleError(w, r, apperrors.NewAppError("invalid occurrence_date, use RFC3339 or YYYY-MM-DD"))
			return
		}
		occurrenceDate = &normalized
	}

	if err := h.srvc.MarkTaskDone(r.Context(), task.ID(id), occurrenceDate); err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Ok(w, map[string]string{"status": "done"})
}

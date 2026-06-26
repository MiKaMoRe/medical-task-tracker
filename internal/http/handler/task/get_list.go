package task

import (
	"net/http"
	"strings"
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

const (
	dateOnlyLayout = "2006-01-02"
)

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if !h.ensureMethod(w, r, http.MethodGet) {
		return
	}

	from, err := parsePeriodDate(r.URL.Query().Get("from"), false)
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid from query param, use RFC3339 or YYYY-MM-DD"))
		return
	}

	to, err := parsePeriodDate(r.URL.Query().Get("to"), true)
	if err != nil {
		h.handleError(w, r, apperrors.NewAppError("invalid to query param, use RFC3339 or YYYY-MM-DD"))
		return
	}

	statuses, err := parseStatusFilters(r.URL.Query()["status"])
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	tasks, err := h.srvc.GetTasks(r.Context(), from, to, statuses)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Ok(w, tasks)
}

func parsePeriodDate(value string, isUpperBound bool) (time.Time, error) {
	if value == "" {
		return time.Time{}, apperrors.NewAppError("query param is required")
	}

	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.UTC(), nil
	}

	if t, err := time.Parse(dateOnlyLayout, value); err == nil {
		if isUpperBound {
			return t.Add(24*time.Hour - time.Nanosecond).UTC(), nil
		}
		return t.UTC(), nil
	}

	return time.Time{}, apperrors.NewAppError("invalid date format")
}

func parseStatusFilters(values []string) ([]domaintask.TaskStatus, error) {
	if len(values) == 0 {
		return nil, nil
	}

	allowedStatuses := map[domaintask.TaskStatus]struct{}{
		domaintask.TaskStatusPlanned: {},
		domaintask.TaskStatusDone:    {},
		domaintask.TaskStatusExpired: {},
	}

	seen := make(map[domaintask.TaskStatus]struct{})
	result := make([]domaintask.TaskStatus, 0, len(values))
	for _, raw := range values {
		for _, part := range strings.Split(raw, ",") {
			normalized := domaintask.TaskStatus(strings.TrimSpace(part))
			if normalized == "" {
				return nil, apperrors.NewAppError("invalid status query param, allowed values: planned, done, expired")
			}
			if _, ok := allowedStatuses[normalized]; !ok {
				return nil, apperrors.NewAppError("invalid status query param, allowed values: planned, done, expired")
			}
			if _, exists := seen[normalized]; exists {
				continue
			}
			seen[normalized] = struct{}{}
			result = append(result, normalized)
		}
	}

	return result, nil
}

func containsTaskStatus(statuses []domaintask.TaskStatus, status domaintask.TaskStatus) bool {
	for _, allowed := range statuses {
		if allowed == status {
			return true
		}
	}
	return false
}

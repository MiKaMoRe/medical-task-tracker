package task

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPut:
		h.UpdateTask(w, r)
	case http.MethodDelete:
		h.DeleteTask(w, r)
	default:
		h.logger.Warn(
			"Method not allowed",
			"method", r.Method,
			"path", r.URL.Path,
			"requestId", requestIDFromRequest(r),
		)
		_ = response.Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}

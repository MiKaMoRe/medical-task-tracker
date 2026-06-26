package task

import (
	"net/http"

	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) ensureMethod(w http.ResponseWriter, r *http.Request, expected string) bool {
	if r.Method == expected {
		return true
	}
	h.logger.Warn(
		"Method not allowed",
		"method", r.Method,
		"path", r.URL.Path,
		"expectedMethod", expected,
		"requestId", requestIDFromRequest(r),
	)
	_ = response.Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	return false
}

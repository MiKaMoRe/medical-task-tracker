package task

import (
	"errors"
	"net/http"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"
)

func (h *TaskHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		valErrs     *apperrors.ValidationMap
		appErr      *apperrors.AppError
		notFoundErr *apperrors.NotFoundError
		conflictErr *apperrors.ConflictError
	)

	reqID := requestIDFromRequest(r)

	switch {
	case errors.As(err, &notFoundErr):
		h.logger.Warn("Task not found", "error", notFoundErr.Message, "requestId", reqID)
		if writeErr := response.NotFound(w, notFoundErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &conflictErr):
		h.logger.Warn("Request conflict", "error", conflictErr.Message, "requestId", reqID)
		if writeErr := response.Error(w, http.StatusConflict, conflictErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &valErrs):
		h.logger.Warn("Validation failed", "error", valErrs.Err().Error(), "requestId", reqID)
		if writeErr := response.UnprocessableEntity(w, valErrs); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &appErr):
		h.logger.Warn("Request rejected", "error", appErr.Message, "requestId", reqID)
		if writeErr := response.UnprocessableEntity(w, appErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	default:
		h.logger.Error("Unhandled error", "error", err.Error(), "requestId", reqID)
		_ = response.InternalServerError(w)
	}

	h.logger.Debug("Request error", "error", err.Error(), "requestId", reqID)
}

func requestIDFromRequest(r *http.Request) string {
	if reqID, ok := httpctx.RequestID(r.Context()); ok {
		return reqID
	}
	return "unknown"
}

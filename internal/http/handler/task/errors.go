package task

import (
	"errors"
	"net/http"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
	"github.com/MiKaMoRe/medical-task-tracker/internal/httpctx"
)

func (h *TaskHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var valErrs *apperrors.ValidationMap
	var appErr *apperrors.AppError
	var notFoundErr *apperrors.NotFoundError
	var conflictErr *apperrors.ConflictError

	reqID, ok := httpctx.RequestID(r.Context())
	if !ok {
		reqID = "unknown"
	}

	switch {
	case errors.As(err, &notFoundErr):
		if writeErr := response.NotFound(w, notFoundErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &conflictErr):
		if writeErr := response.Error(w, http.StatusConflict, conflictErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &valErrs):
		if writeErr := response.UnprocessableEntity(w, valErrs); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	case errors.As(err, &appErr):
		if writeErr := response.UnprocessableEntity(w, appErr.Message); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error(), "requestId", reqID)
		}
	default:
		h.logger.Error("Unhandled error", "error", err.Error(), "requestId", reqID)
		_ = response.InternalServerError(w)
	}

	h.logger.Debug("Request error", "error", err.Error(), "requestId", reqID)
}

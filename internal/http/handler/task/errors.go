package task

import (
	"errors"
	"net/http"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) handleError(w http.ResponseWriter, err error) {
	var valErrs *apperrors.ValidationErrors
	var appErr *apperrors.AppError
	var notFoundErr *apperrors.NotFoundError
	var conflictErr *apperrors.ConflictError

	switch {
	case errors.As(err, &notFoundErr):
		msg := map[string]string{"error": notFoundErr.Message}
		if writeErr := response.NotFound(w, msg); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error())
		}
	case errors.As(err, &conflictErr):
		msg := map[string]string{"error": conflictErr.Message}
		if writeErr := response.Error(w, http.StatusConflict, msg); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error())
		}
	case errors.As(err, &valErrs):
		msg := map[string][]apperrors.ValidationError{"validation_errors": valErrs.Errors}
		if writeErr := response.UnprocessableEntity(w, msg); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error())
		}
	case errors.As(err, &appErr):
		msg := map[string]string{"error": appErr.Message}
		if writeErr := response.UnprocessableEntity(w, msg); writeErr != nil {
			h.logger.Error("Response error", "error", writeErr.Error())
		}
	default:
		h.logger.Error("Unhandled error", "error", err.Error())
		_ = response.InternalServerError(w)
	}

	h.logger.Debug("Request error", "error", err.Error())
}

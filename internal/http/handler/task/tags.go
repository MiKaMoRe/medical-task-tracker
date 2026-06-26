package task

import (
	"encoding/json"
	"net/http"
	"strconv"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	domaintask "github.com/MiKaMoRe/medical-task-tracker/internal/domain/task"
	"github.com/MiKaMoRe/medical-task-tracker/internal/http/response"
)

func (h *TaskHandler) TaskTags(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AddTaskTags(w, r)
	case http.MethodDelete:
		h.RemoveTaskTags(w, r)
	default:
		_ = response.Error(w, http.StatusMethodNotAllowed, "Method Not Allowed")
	}
}

func (h *TaskHandler) AddTaskTags(w http.ResponseWriter, r *http.Request) {
	id, req, err := decodeTagsRequest(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	updated, err := h.srvc.AddTaskTags(r.Context(), id, req.Tags)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Ok(w, updated)
}

func (h *TaskHandler) RemoveTaskTags(w http.ResponseWriter, r *http.Request) {
	id, req, err := decodeTagsRequest(r)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	updated, err := h.srvc.RemoveTaskTags(r.Context(), id, req.Tags)
	if err != nil {
		h.handleError(w, r, err)
		return
	}

	response.Ok(w, updated)
}

func decodeTagsRequest(r *http.Request) (domaintask.ID, TaskTagsRequest, error) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return domaintask.IDFromInt(0), TaskTagsRequest{}, apperrors.NewAppError("invalid task id")
	}

	var req TaskTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return domaintask.IDFromInt(0), TaskTagsRequest{}, apperrors.NewAppError("invalid request body")
	}

	return domaintask.IDFromInt(id), req, nil
}

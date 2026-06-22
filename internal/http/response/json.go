package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, status int, data any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return fmt.Errorf("response: encode JSON: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := buf.WriteTo(w); err != nil {
		return fmt.Errorf("response: write body: %w", err)
	}

	return nil
}

func Error(w http.ResponseWriter, status int, message any) error {
	return JSON(w, status, map[string]any{"error": message})
}

func Ok(w http.ResponseWriter, data any) error {
	return JSON(w, http.StatusOK, data)
}

func Created(w http.ResponseWriter, data any) error {
	return JSON(w, http.StatusCreated, data)
}

func BadRequest(w http.ResponseWriter) error {
	return Error(w, http.StatusBadRequest, "Bad Request")
}

func NotFound(w http.ResponseWriter, message any) error {
	return Error(w, http.StatusNotFound, message)
}

func InternalServerError(w http.ResponseWriter) error {
	return Error(w, http.StatusInternalServerError, "Internal Server Error")
}

func UnprocessableEntity(w http.ResponseWriter, message any) error {
	return Error(w, http.StatusUnprocessableEntity, message)
}

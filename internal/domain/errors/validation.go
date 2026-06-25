package apperrors

import "fmt"

type ValidationErrors struct {
	Errors []ValidationError
}

type ValidationError struct {
	Field   string
	Message string
}

func Validation(errs ...ValidationError) error {
	return &ValidationErrors{Errors: errs}
}

func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %d error(s)", len(e.Errors))
}

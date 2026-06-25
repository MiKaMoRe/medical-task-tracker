package apperrors

type ConflictError struct {
	Message string
}

func (e *ConflictError) Error() string {
	return e.Message
}

func Conflict(message string) error {
	return &ConflictError{Message: message}
}

package apperrors

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func NotFound(message string) error {
	return &NotFoundError{Message: message}
}

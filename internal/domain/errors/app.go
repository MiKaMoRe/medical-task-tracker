package apperrors

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(message string) error {
	return &AppError{message}
}

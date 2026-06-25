package vo

import (
	"fmt"
	"strings"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
)

func TrimmedText(raw, field string, minLen, maxLen int) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) < minLen || len(trimmed) > maxLen {
		return "", apperrors.Validation(apperrors.ValidationError{
			Field:   field,
			Message: fmt.Sprintf("must be between %d and %d characters", minLen, maxLen),
		})
	}
	return trimmed, nil
}

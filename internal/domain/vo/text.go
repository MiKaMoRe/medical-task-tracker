package vo

import (
	"fmt"
	"strings"
)

func TrimmedText(raw, field string, minLen, maxLen int) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) < minLen || len(trimmed) > maxLen {
		return "", fmt.Errorf("must be between %d and %d characters", minLen, maxLen)
	}
	return trimmed, nil
}

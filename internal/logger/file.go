package logger

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultLogFile = "logs/dev.log"

func logFilePath() string {
	if p := os.Getenv("LOG_FILE_PATH"); p != "" {
		return p
	}
	return defaultLogFile
}

func openLogFile() (*os.File, error) {
	path := logFilePath()
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("logger: create log dir %q: %w", dir, err)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("logger: open log file %q: %w", path, err)
	}

	return f, nil
}

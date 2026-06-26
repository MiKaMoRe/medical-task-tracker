package task

import (
	"testing"
	"time"
)

func TestResolveTaskStatus_ReturnsExpiredForPastDay(t *testing.T) {
	now := time.Date(2026, time.June, 26, 15, 0, 0, 0, time.UTC)
	taskDate := time.Date(2026, time.June, 25, 23, 59, 59, 0, time.UTC)

	status := ResolveTaskStatus(taskDate, false, now)

	if status != TaskStatusExpired {
		t.Fatalf("expected status %q, got %q", TaskStatusExpired, status)
	}
}

func TestResolveTaskStatus_DoneHasPriorityOverExpired(t *testing.T) {
	now := time.Date(2026, time.June, 26, 15, 0, 0, 0, time.UTC)
	taskDate := time.Date(2026, time.June, 20, 8, 0, 0, 0, time.UTC)

	status := ResolveTaskStatus(taskDate, true, now)

	if status != TaskStatusDone {
		t.Fatalf("expected status %q, got %q", TaskStatusDone, status)
	}
}

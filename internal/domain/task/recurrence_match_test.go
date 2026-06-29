package task

import (
	"testing"
	"time"
)

func TestRecurrenceRuleMatches(t *testing.T) {
	t.Run("weekly", func(t *testing.T) {
		rule := WeeklyRecurrenceRule{WeekDay: 1}
		if !rule.Matches(mustUTC(2026, time.June, 1, 9, 0, 0), time.Time{}) {
			t.Fatalf("expected weekly match")
		}
		if rule.Matches(mustUTC(2026, time.June, 2, 9, 0, 0), time.Time{}) {
			t.Fatalf("expected weekly mismatch")
		}
	})

	t.Run("shift", func(t *testing.T) {
		rule := ShiftRecurrenceRule{NumberOfTaskDays: 2, NumberOfShiftDays: 2}
		start := mustUTC(2026, time.June, 1, 9, 0, 0)
		if !rule.Matches(mustUTC(2026, time.June, 2, 9, 0, 0), start) {
			t.Fatalf("expected shift day to match")
		}
		if rule.Matches(mustUTC(2026, time.June, 3, 9, 0, 0), start) {
			t.Fatalf("expected shift day to be off-cycle")
		}
	})
}

func mustUTC(year int, month time.Month, day int, hour int, minute int, second int) time.Time {
	return time.Date(year, month, day, hour, minute, second, 0, time.UTC)
}

package task

import (
	"strings"
	"testing"
	"time"
)

func TestNewDate_UsesUTCCalendarDayComparison(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	todayStartUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayUTC := todayStartUTC.Add(-24 * time.Hour)

	tests := []struct {
		name          string
		input         time.Time
		wantErrSubstr string
	}{
		{
			name:  "accepts date earlier today in utc",
			input: todayStartUTC,
		},
		{
			name:          "rejects previous utc day",
			input:         yesterdayUTC,
			wantErrSubstr: "date cannot be in the past",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, errs := NewDate(tc.input)
			if tc.wantErrSubstr == "" {
				if len(errs) != 0 {
					t.Fatalf("expected no errors, got %v", errs)
				}
				return
			}

			if !containsError(errs, tc.wantErrSubstr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErrSubstr, errs)
			}
		})
	}
}

func TestNewRecurringTask_EndDateUsesUTCCalendarDayComparison(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	todayStartUTC := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterdayUTC := todayStartUTC.Add(-24 * time.Hour)
	weekDay := int(todayStartUTC.Weekday())

	tests := []struct {
		name          string
		endDate       time.Time
		wantErrSubstr string
	}{
		{
			name:    "accepts end date earlier today in utc",
			endDate: todayStartUTC,
		},
		{
			name:          "rejects previous utc day",
			endDate:       yesterdayUTC,
			wantErrSubstr: "end date must be after current date",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, vm := NewRecurringTask(
				string(RecurringTypeWeekly),
				&tc.endDate,
				RecurrenceRuleInput{
					Weekly: &WeeklyRecurrenceRuleInput{
						WeekDay: &weekDay,
					},
				},
			)
			err := vm.Err()
			if tc.wantErrSubstr == "" {
				if err != nil {
					t.Fatalf("expected no errors, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErrSubstr) {
				t.Fatalf("expected error containing %q, got %v", tc.wantErrSubstr, err)
			}
		})
	}
}

func containsError(errs []error, want string) bool {
	for _, err := range errs {
		if err != nil && strings.Contains(err.Error(), want) {
			return true
		}
	}
	return false
}

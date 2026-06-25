package task

import (
	"errors"
	"time"
)

type RecurringType string

const (
	RecurringTypeWeekly   RecurringType = "weekly"
	RecurringTypeMonthly  RecurringType = "monthly"
	RecurringTypeYearly   RecurringType = "yearly"
	RecurringTypeBiweekly RecurringType = "biweekly"
	RecurringTypeShift    RecurringType = "shift"
	RecurringTypeParity   RecurringType = "parity"
)

func (r RecurringType) String() string {
	return string(r)
}

func NewRecurringType(s string) (RecurringType, error) {
	switch s {
	case "weekly":
		return RecurringTypeWeekly, nil
	case "monthly":
		return RecurringTypeMonthly, nil
	case "yearly":
		return RecurringTypeYearly, nil
	case "biweekly":
		return RecurringTypeBiweekly, nil
	case "shift":
		return RecurringTypeShift, nil
	case "parity":
		return RecurringTypeParity, nil
	default:
		return RecurringType(s), errors.New("invalid recurring type")
	}
}

type RecurringTask struct {
	RecurringType RecurringType `json:"recurring_type"`
	EndDate       Date          `json:"end_date"`
}

func NewRecurringTask(recurringType string, startDateString string, endDateString string) (*RecurringTask, error) {
	startDate, err := NewDate(startDateString)
	if err != nil {
		return nil, err
	}

	oneDayAgo := Date(time.Now().In(time.UTC).Add(-24 * time.Hour))
	if startDate.IsAfter(oneDayAgo) {
		return nil, errors.New("start date must be before current date minus one day")
	}

	endDate, err := NewDate(endDateString)
	if err != nil {
		return nil, err
	}

	if endDate.IsBefore(Date(time.Now().In(time.UTC))) {
		return nil, errors.New("end date must be after current date")
	}

	if endDate.IsBefore(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	rType, err := NewRecurringType(recurringType)
	if err != nil {
		return nil, err
	}

	return &RecurringTask{
		RecurringType: rType,
		EndDate:       endDate,
	}, nil
}

package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Date time.Time

func (d Date) IsBefore(other Date) bool {
	return time.Time(d).Before(time.Time(other))
}

func (d Date) IsAfter(other Date) bool {
	return time.Time(d).After(time.Time(other))
}

func NewDate(datetime time.Time) (Date, []error) {
	errs := []error{}
	if datetime.IsZero() {
		errs = append(errs, errors.New("date cannot be zero"))
	}

	date := Date(datetime.In(time.UTC))
	if floorDate(time.Time(date)).Before(floorDate(time.Now().In(time.UTC))) {
		errs = append(errs, errors.New("date cannot be in the past"))
	}

	return date, errs
}

func DateFromDB(date string) (Date, error) {
	dateTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Date(time.Time{}), fmt.Errorf("invalid date: %w", err)
	}

	return Date(dateTime), nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d))
}

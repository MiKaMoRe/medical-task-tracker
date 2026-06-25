package task

import "time"

type Date time.Time

func (d Date) IsBefore(other Date) bool {
	return d.IsBefore(other)
}

func (d Date) IsAfter(other Date) bool {
	return d.IsAfter(other)
}

func NewDate(date string) (Date, error) {
	dateTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Date(time.Time{}), err
	}
	return Date(dateTime.In(time.UTC)), nil
}

func DateFromDB(date string) (Date, error) {
	dateTime, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return Date(time.Time{}), err
	}
	return Date(dateTime), nil
}

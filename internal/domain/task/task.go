package task

import (
	"errors"

	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
)

type Task struct {
	ID          ID             `json:"id"`
	Title       Title          `json:"title"`
	Description Description    `json:"description"`
	Date        Date           `json:"date"`
	IsRecurring bool           `json:"is_recurring"`
	Recurring   *RecurringTask `json:"recurring,omitempty"`
	Tags        []tag.Tag      `json:"tags"`
}

func NewTask(
	title Title,
	description Description,
	date Date,
	dueDate Date,
	isRecurring bool,
	recurring *RecurringTask,
	tags []string,
) (*Task, error) {
	if date.IsAfter(dueDate) {
		return nil, errors.New("date must be before due date")
	}

	ntags, err := tag.NewTags(tags)
	if err != nil {
		return nil, err
	}

	return &Task{
		Title:       title,
		Description: description,
		Date:        date,
		IsRecurring: isRecurring,
		Tags:        ntags,
	}, nil
}

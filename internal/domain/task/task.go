package task

import (
	"time"

	apperrors "github.com/MiKaMoRe/medical-task-tracker/internal/domain/errors"
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/tag"
)

type Task struct {
	ID          ID                  `json:"id"`
	Title       Title               `json:"title"`
	Description Description         `json:"description"`
	Date        Date                `json:"date"`
	IsRecurring bool                `json:"is_recurring"`
	Status      TaskStatus          `json:"status"`
	Recurring   *RecurringTask      `json:"recurring,omitempty"`
	Tags        []tag.Tag           `json:"tags"`
	IsDone      bool                `json:"-"`
	DoneDates   map[string]struct{} `json:"-"`
}

func NewTask(
	title string,
	description string,
	date time.Time,
	isRecurring bool,
	tags []string,
) (*Task, *apperrors.ValidationMap) {
	vm := apperrors.NewValidationMap()
	nTitle, err := NewTitle(title)
	vm.Add("task.title", err)

	nDescription, err := NewDescription(description)
	vm.Add("task.description", err)

	nDate, errs := NewDate(date)
	vm.Add("task.date", errs...)

	nTags, errs := tag.NewTags(tags)
	vm.Add("task.tags", errs...)

	return &Task{
		Title:       nTitle,
		Description: nDescription,
		Date:        nDate,
		IsRecurring: isRecurring,
		Tags:        nTags,
	}, vm
}

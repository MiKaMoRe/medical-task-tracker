package task

import (
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/vo"
)

type Description string

const descriptionMinLen = 1
const descriptionMaxLen = 1000

func NewDescription(description string) (Description, error) {
	trimmed, err := vo.TrimmedText(description, "description", descriptionMinLen, descriptionMaxLen)
	if err != nil {
		return "", err
	}
	return Description(trimmed), nil
}

func DescriptionFromDB(description string) (Description, error) {
	return Description(description), nil
}

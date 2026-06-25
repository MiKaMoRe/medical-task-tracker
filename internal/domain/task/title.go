package task

import "github.com/MiKaMoRe/medical-task-tracker/internal/domain/vo"

type Title string

const titleMinLen = 1
const titleMaxLen = 200

func NewTitle(title string) (Title, error) {
	trimmed, err := vo.TrimmedText(title, "title", titleMinLen, titleMaxLen)
	if err != nil {
		return "", err
	}
	return Title(trimmed), nil
}

func NameFromDB(title string) (Title, error) {
	return Title(title), nil
}

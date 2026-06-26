package tag

import (
	"github.com/MiKaMoRe/medical-task-tracker/internal/domain/vo"
)

type Name string

func NewName(name string) (Name, error) {
	trimmed, err := vo.TrimmedText(name, "name", 1, 255)
	if err != nil {
		return "", err
	}
	return Name(trimmed), nil
}

func (n Name) String() string {
	return string(n)
}

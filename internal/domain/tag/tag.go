package tag

import (
	"errors"
)

type Tag struct {
	ID   ID   `json:"id"`
	Name Name `json:"name"`
}

func NewTag(name string) (Tag, []error) {
	errs := []error{}
	if len(name) > 255 {
		errs = append(errs, errors.New("name must be less than 255 characters"))
	}

	if name == "" {
		errs = append(errs, errors.New("name cannot be empty"))
	}

	nName, err := NewName(name)
	if err != nil {
		errs = append(errs, err)
	}
	return Tag{Name: nName}, errs
}

func NewTags(names []string) ([]Tag, []error) {
	tags := make([]Tag, len(names))
	errs := make([]error, len(names))
	for i, name := range names {
		tag, err := NewTag(name)
		if err != nil {
			errs = append(errs, err...)
		}
		tags[i] = tag
	}
	return tags, errs
}

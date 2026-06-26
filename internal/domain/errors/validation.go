package apperrors

import (
	"fmt"
	"strings"
)

type ValidationMap struct {
	Errors map[string][]string `json:"validation_errors"`
}

func NewValidationMap() *ValidationMap {
	return &ValidationMap{Errors: make(map[string][]string)}
}

func (e *ValidationMap) Add(field string, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}
		e.Errors[field] = append(e.Errors[field], err.Error())
	}
}

func MergeValidationMaps(maps ...*ValidationMap) *ValidationMap {
	vm := NewValidationMap()
	for _, m := range maps {
		for field, msgs := range m.Errors {
			vm.Errors[field] = append(vm.Errors[field], msgs...)
		}
	}
	return vm
}

func (e *ValidationMap) Error() string {
	errors := make([]string, 0, len(e.Errors))
	for field, errs := range e.Errors {
		for _, err := range errs {
			errors = append(errors, fmt.Sprintf("%s: %s", field, err))
		}
	}
	return strings.Join(errors, "\n")
}

func (e *ValidationMap) IsEmpty() bool {
	return len(e.Errors) == 0
}

func (e *ValidationMap) Err() error {
	if e.IsEmpty() {
		return nil
	}
	return e
}

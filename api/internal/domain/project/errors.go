package project

import "errors"

var ErrProjectNotFound = errors.New("project not found")

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string { return "validation failed" }

func (e *ValidationError) add(field, message string) {
	if e.Fields == nil {
		e.Fields = map[string]string{}
	}
	if _, exists := e.Fields[field]; !exists {
		e.Fields[field] = message
	}
}

func (e *ValidationError) orNil() error {
	if len(e.Fields) == 0 {
		return nil
	}
	return e
}

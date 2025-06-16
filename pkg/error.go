package config

import (
	"errors"
	"strings"
)

// Enumeration of errors that may be returned by configuration operations.
const (
	ErrNotInitialized = configErr("not initialized")
	ErrPathNotFound   = configErr("path not found")
)

// configErr defines the type for errors that may be returned by configuration operations.
type configErr string

// Error returns the cause of the file system error.
func (e configErr) Error() string {
	return string(e)
}

// PathError is used for recording errors that may occur when parsing values from configuration paths.
type PathError struct {
	Err       error
	Operation string
	Path      string
}

// Error returns the error message for the PathError.
func (e *PathError) Error() string {
	var pe strings.Builder
	pe.WriteString("configuration: ")

	if op := strings.TrimSpace(e.Operation); op != "" {
		pe.WriteString(e.Operation + ": ")
	}

	if e.Path != "" {
		pe.WriteString(e.Path + ": ")
	}

	if e.Err == nil {
		e.Err = errors.New("invalid")
	}
	pe.WriteString(e.Err.Error())
	return pe.String()
}

func (e *PathError) Unwrap() error {
	return e.Err
}

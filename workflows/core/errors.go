// Package core provides error types for workflow operations.
package core

import "fmt"

// LoadError is returned when loading a workflow repository fails.
type LoadError struct {
	Path string
	Err  error
}

func (e *LoadError) Error() string {
	return fmt.Sprintf("failed to load workflow repo from %s: %v", e.Path, e.Err)
}

func (e *LoadError) Unwrap() error {
	return e.Err
}

// NotFoundError is returned when a workflow is not found.
type NotFoundError struct {
	WorkflowID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("workflow not found: %s", e.WorkflowID)
}

// AdapterNotFoundError is returned when an adapter is not found.
type AdapterNotFoundError struct {
	Name string
}

func (e *AdapterNotFoundError) Error() string {
	return fmt.Sprintf("workflow adapter not found: %s", e.Name)
}

// GenerateError is returned when generation fails.
type GenerateError struct {
	Platform string
	Path     string
	Err      error
}

func (e *GenerateError) Error() string {
	return fmt.Sprintf("failed to generate %s workflow at %s: %v", e.Platform, e.Path, e.Err)
}

func (e *GenerateError) Unwrap() error {
	return e.Err
}

// ParseError is returned when parsing fails.
type ParseError struct {
	Path string
	Err  error
}

func (e *ParseError) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("failed to parse workflow at %s: %v", e.Path, e.Err)
	}
	return fmt.Sprintf("failed to parse workflow: %v", e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// WriteError is returned when writing fails.
type WriteError struct {
	Path string
	Err  error
}

func (e *WriteError) Error() string {
	return fmt.Sprintf("failed to write workflow to %s: %v", e.Path, e.Err)
}

func (e *WriteError) Unwrap() error {
	return e.Err
}

package task

import (
	"errors"

	"github.com/LNMMusic/optional"
)

// Interfaces
type Task struct {
	ID 			optional.Option[string]
	Title 		optional.Option[string]
	Description optional.Option[string]
	Completed 	optional.Option[bool]
}

// Storage is the interface that wraps the basic methods for a task storage.
type Storage interface {
	// Get returns the task with the given id.
	Get(id string) (ts *Task, err error)

	// Save saves the given task.
	Save(task *Task) (err error)
}
var (
	ErrStorageInternal = errors.New("storage internal error")
	ErrStorageNotFound = errors.New("storage task not found")
	ErrStorageInvalid  = errors.New("storage invalid task")
)

// Validator is the interface that wraps the basic methods for a task validator.
type Validator interface {
	// Validate validates the given task.
	Validate(task *Task) (err error)
}
var (
	ErrValidatorInternal 	  = errors.New("validator internal error")
	ErrValidatorFieldRequired = errors.New("validator field required")
	ErrValidatorFieldEmpty	  = errors.New("validator field empty")
	ErrValidatorFieldQuality  = errors.New("validator field quality")
)
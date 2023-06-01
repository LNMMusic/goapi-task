package task

import (
	"fmt"

	"github.com/LNMMusic/optional"

	"github.com/google/uuid"
)

// constructor
func NewStorageLocal(db []*Task, vl Validator) *StorageLocal {
	return &StorageLocal{db: db, vl: vl}
}


// StorageLocal is the local implementation of the task storage.
type StorageLocal struct {
	db []*Task
	vl Validator
}

func (s *StorageLocal) Get(id string) (ts *Task, err error) {
	for _, t := range s.db {
		tId, _ := t.ID.Unwrap()
		if tId == id {
			ts = t
			return
		}
	}

	err = fmt.Errorf("%w: %v", ErrStorageNotFound, id)
	return
}

func (s *StorageLocal) Save(task *Task) (err error) {
	// validate task
	err = s.vl.Validate(task)
	if err != nil {
		err  = fmt.Errorf("%w: %v", ErrStorageInvalid, err)
		return
	}

	// generate id
	id := uuid.New().String()
	task.ID = optional.Some(id)

	// save task
	s.db = append(s.db, task)
	return
}
	
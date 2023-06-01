package task

import "fmt"

// constructor
func NewValidatorLocal() *ValidatorLocal {
	return &ValidatorLocal{}
}

// ValidatorLocal is the local implementation of the task validator.
type ValidatorLocal struct {}

func (v *ValidatorLocal) Validate(task *Task) (err error) {
	// check required fields (non nullable)
	if !task.Title.IsSome() {
		err = fmt.Errorf("%w: title", ErrValidatorFieldRequired)
		return
	}
	if !task.Completed.IsSome() {
		err = fmt.Errorf("%w: completed", ErrValidatorFieldRequired)
		return
	}

	// check empty fields and quality values (non nullable)
	// -> safe to not check err, due to the previous check
	title, _ := task.Title.Unwrap()
	if title == "" {
		err = fmt.Errorf("%w: title", ErrValidatorFieldEmpty)
		return
	}
	if len(title) > 50 {
		err = fmt.Errorf("%w: title", ErrValidatorFieldQuality)
		return
	}

	// check empty fields and quality values (nullable)
	// -> safe to not check err, due to the previous check
	if task.Description.IsSome() {
		var description string
		description, _ = task.Description.Unwrap()
		if len(description) > 150 {
			err = fmt.Errorf("%w: description", ErrValidatorFieldQuality)
			return
		}
	}

	return
}
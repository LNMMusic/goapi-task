package task

import "github.com/stretchr/testify/mock"

// constructor
func NewValidatorMock() *ValidatorMock {
	return &ValidatorMock{}
}

// ValidatorMock is the mock implementation of the task validator.
type ValidatorMock struct {
	mock.Mock
}

func (v *ValidatorMock) Validate(task *Task) (err error) {
	args := v.Called(task)
	return args.Error(0)
}
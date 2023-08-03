package profiles

import "github.com/stretchr/testify/mock"

// NewImplValidatorMock returns a new mock for the Validator interface
func NewImplValidatorMock() *ImplValidatorMock {
	return &ImplValidatorMock{}
}

// ImplValidatorMock is a mock implementation of the Validator interface
type ImplValidatorMock struct {
	mock.Mock
}

// Default set default values for a profile
func (mk *ImplValidatorMock) Default(pf *Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}

// Validate validates a profile
func (mk *ImplValidatorMock) Validate(pf *Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}
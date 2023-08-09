package validator

import (
	"api/internal/profiles"

	"github.com/stretchr/testify/mock"
)

// NewImplProfilesValidatorMock returns a new mock for the Validator interface
func NewImplProfilesValidatorMock() *ImplProfilesValidatorMock {
	return &ImplProfilesValidatorMock{}
}

// ImplProfilesValidatorMock is a mock implementation of the Validator interface
type ImplProfilesValidatorMock struct {
	mock.Mock
}

// Default set default values for a profile
func (mk *ImplProfilesValidatorMock) Default(pf *profiles.Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}

// Validate validates a profile
func (mk *ImplProfilesValidatorMock) Validate(pf *profiles.Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}
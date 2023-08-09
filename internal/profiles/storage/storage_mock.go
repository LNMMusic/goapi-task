package storage

import (
	"api/internal/profiles"

	"github.com/stretchr/testify/mock"
)

// NewImplProfilesStorageMock returns a new mock for the Storage interface
func NewImplProfilesStorageMock() *ImplProfilesStorageMock {
	return &ImplProfilesStorageMock{}
}

// ImplProfilesStorageMock is a mock implementation of the Storage interface
type ImplProfilesStorageMock struct {
	mock.Mock
}

// GetProfileByUserId provides a mock function with given fields: userId
func (mk *ImplProfilesStorageMock) GetProfileById(id string) (pf *profiles.Profile, err error) {
	args := mk.Called(id)
	pf = args.Get(0).(*profiles.Profile)
	err = args.Error(1)
	return
}

// ActivateProfile provides a mock function with given fields: pf
func (mk *ImplProfilesStorageMock) ActivateProfile(pf *profiles.Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}
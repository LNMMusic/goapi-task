package profiles

import "github.com/stretchr/testify/mock"

// NewImplStorageMock returns a new mock for the Storage interface
func NewImplStorageMock() *ImplStorageMock {
	return &ImplStorageMock{}
}

// ImplStorageMock is a mock implementation of the Storage interface
type ImplStorageMock struct {
	mock.Mock
}

// GetProfileByUserId provides a mock function with given fields: userId
func (mk *ImplStorageMock) GetProfileByUserId(userId string) (pf *Profile, err error) {
	args := mk.Called(userId)
	pf = args.Get(0).(*Profile)
	err = args.Error(1)
	return
}

// ActivateProfile provides a mock function with given fields: pf
func (mk *ImplStorageMock) ActivateProfile(pf *Profile) (err error) {
	args := mk.Called(pf)
	err = args.Error(0)
	return
}
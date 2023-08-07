package uuidgenerator

import "github.com/stretchr/testify/mock"

// NewUUIDGeneratorMock returns a new UUIDGeneratorMock
func NewUUIDGeneratorMock() (ug *ImplUUIDGeneratorMock) {
	ug = &ImplUUIDGeneratorMock{}
	return
}

// ImplUUIDGeneratorMock is the implementation of UUIDGenerator using a mock UUID generator
type ImplUUIDGeneratorMock struct{
	mock.Mock
}

func (ug *ImplUUIDGeneratorMock) UUID() (id string) {
	args := ug.Called()
	id = args.String(0)
	return
}
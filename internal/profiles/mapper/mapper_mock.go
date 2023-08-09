package mapper

import "github.com/stretchr/testify/mock"

// NewProfileMapperMock returns a new ProfileMapperMock
func NewProfileMapperMock() *ProfileMapperMock {
	return &ProfileMapperMock{}
}

// ProfileMapperMock is the mock for ProfileMapper
type ProfileMapperMock struct {
	mock.Mock
}

// MapProfile maps a profile
func (m *ProfileMapperMock) MapProfile(userId string) (profileId string, err error) {
	args := m.Called(userId)
	profileId = args.String(0)
	err = args.Error(1)
	return
}
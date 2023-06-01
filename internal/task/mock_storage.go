package task

import "github.com/stretchr/testify/mock"

// constructor
func NewStorageMock() *StorageMock {
	mk := &StorageMock{}
	mk.SetTask = func(t *Task) {}
	return mk
}

// StorageMock is a mock implementation of the task storage.
type StorageMock struct {
	mock.Mock
	SetTask func(t *Task)
}

func (m *StorageMock) Get(id string) (ts *Task, err error) {
	args := m.Called(id)
	ts = args.Get(0).(*Task)
	err = args.Error(1)
	return
}

func (m *StorageMock) Save(t *Task) (err error) {
	args := m.Called(t)

	m.SetTask(t)

	err = args.Error(0)
	return
}
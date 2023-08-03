package transactioner

import "github.com/stretchr/testify/mock"

// NewImplTransactionerMock returns a new mock for the Transactioner interface
func NewImplTransactionerMock() *ImplTransactionerMock {
	return &ImplTransactionerMock{}
}

// ImplTransactionerMock is a mock implementation of the Transactioner interface
type ImplTransactionerMock struct {
	mock.Mock
}

// Do provides a mock function with given fields: f
func (mk *ImplTransactionerMock) Do(op operation) (err error) {
	args := mk.Called(op)

	op()

	err = args.Error(0)
	return
}
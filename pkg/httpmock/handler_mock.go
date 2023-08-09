package httpmock

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// NewHandlerMock returns a new HandlerMock
func NewHandlerMock() *HandlerMock {
	return &HandlerMock{
		SetUpServeHTTP: func(w http.ResponseWriter, r *http.Request) {},
	}
}

// HandlerMock is a mock for http.Handler
type HandlerMock struct {
	mock.Mock
	SetUpServeHTTP func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP is a mock for http.Handler.ServeHTTP
func (m *HandlerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
	m.SetUpServeHTTP(w, r)
}
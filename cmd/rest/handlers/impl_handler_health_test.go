package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests
func TestHandler_Health(t *testing.T) {
	t.Run("Health", func(t *testing.T) {
		// arrange
		// ...

		// act
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/ping", nil)
		Health()(w, r)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "pong", w.Body.String())			
	})
}
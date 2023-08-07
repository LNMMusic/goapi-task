package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for JSON functionality
func TestJSON(t *testing.T) {
	type body struct { Message string; Data string; Error bool}
	type input struct { w *httptest.ResponseRecorder; code int; body body }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
	}

	cases := []testCase{
		// case: 200 OK
		{
			name: "200 OK",
			input: input{
				w: httptest.NewRecorder(),
				code: http.StatusOK,	
				body: body{
					Message: "Success",
					Data:    "data",
					Error:   false,
				},
			},
			output: output{
				code: http.StatusOK,
				body: `{"Message":"Success","Data":"data","Error":false}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		},
		// case: 400 Bad Request
		{
			name: "400 Bad Request",
			input: input{
				w: httptest.NewRecorder(),
				code: http.StatusBadRequest,
				body: body{
					Message: "Bad Request",
					Data:    "data",
					Error:   true,
				},
			},
			output: output{
				code: http.StatusBadRequest,
				body: `{"Message":"Bad Request","Data":"data","Error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		},
		// case: 500 Internal Server Error
		{
			name: "500 Internal Server Error",
			input: input{
				w: httptest.NewRecorder(),
				code: http.StatusInternalServerError,
				body: body{
					Message: "Internal Server Error",
					Data:    "data",
					Error:   true,
				},
			},
			output: output{
				code: http.StatusInternalServerError,
				body: `{"Message":"Internal Server Error","Data":"data","Error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			// ...

			// act
			JSON(c.input.w, c.input.code, c.input.body)

			// assert
			assert.Equal(t, c.output.code, c.input.w.Code)
			assert.JSONEq(t, c.output.body, c.input.w.Body.String())
			assert.Equal(t, c.output.headers, c.input.w.Header())
		})
	}
}
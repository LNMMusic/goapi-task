package response

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests
func TestResponse_Ok(t *testing.T) {
	type input struct {code int; msg string; data any}
	type output struct {code int; header string; response string}
	type testCase struct {
		title  string
		input  input
		output output
	}

	cases := []testCase{
		{
			title: "response full",
			input: input{
				code: 200,
				msg:  "ok",
				data: "data",
			},
			output: output{
				code: 200,
				header: "application/json",
				response: `{"message":"ok","data":"data"}`,
			},
		},
		{
			title: "response without data",
			input: input{
				code: 200,
				msg:  "ok",
				data: nil,
			},
			output: output{
				code: 200,
				header: "application/json",
				response: `{"message":"ok","data":null}`,
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			w := httptest.NewRecorder()

			// act
			Ok(w, c.input.code, c.input.msg, c.input.data)

			// assert
			assert.Equal(t, c.output.code, w.Code)
			assert.Equal(t, c.output.header, w.Header().Get("Content-Type"))
			assert.JSONEq(t, c.output.response, w.Body.String())
		})
	}
}

func TestResponse_Err(t *testing.T) {
	type input struct {code int; msg string}
	type output struct {code int; header string; response string}
	type testCase struct {
		title  string
		input  input
		output output
	}

	cases := []testCase{
		{
			title: "response full",
			input: input{
				code: 400,
				msg:  "err",
			},
			output: output{
				code: 400,
				header: "application/json",
				response: `{"data":null,"message":"err"}`,
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			w := httptest.NewRecorder()

			// act
			Err(w, c.input.code, c.input.msg)

			// assert
			assert.Equal(t, c.output.code, w.Code)
			assert.Equal(t, c.output.header, w.Header().Get("Content-Type"))
			assert.JSONEq(t, c.output.response, w.Body.String())
		})
	}
}
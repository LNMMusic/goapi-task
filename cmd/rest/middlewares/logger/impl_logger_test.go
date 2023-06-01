package logger

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests
func TestLog_Display(t *testing.T) {
	type output struct {display string}
	type testCase struct {
		title  string
		input  any
		output output
		setLog func(lg *Log)
	}

	cases := []testCase{
		{
			title: "empty log",
			input: nil,
			output: output{display: "[]  0 | ip:  | errors: []"},
			setLog: func(lg *Log) {},
		},
		{
			title: "log with ip, method GET, path, status, and errors",
			input: nil,
			output: output{display: "[GET] /tasks 200 | ip: 0.0.0.0 | errors: []"},
			setLog: func(lg *Log) {
				lg.IP = "0.0.0.0"
				lg.Method = "GET"
				lg.Path = "/tasks"
				lg.Status = 200
				lg.Errors = []error{}
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			lg := &Log{}
			c.setLog(lg)

			// act
			display := lg.Display()

			// assert
			assert.Equal(t, c.output.display, display)
		})
	}
}

func TestLogger_Errors(t *testing.T) {
	var errs []error = []error{errors.New("error 1"), errors.New("error 2")}

	type input struct {r *http.Request; errs []error}
	type output struct {errs []error}
	type testCase struct {
		title  string
		input  input
		output output
	}

	cases := []testCase{
		{
			title: "empty errors",
			input: input{r: &http.Request{}, errs: []error{}},
			output: output{errs: []error{}},
		},
		{
			title: "1 error",
			input: input{r: &http.Request{}, errs: []error{errs[0]}},
			output: output{errs: []error{errs[0]}},
		},
		{
			title: "+1 errors",
			input: input{r: &http.Request{}, errs: errs},
			output: output{errs: errs},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			// ...

			// act
			Errors(c.input.r, c.input.errs...)

			// assert
			errs, ok := c.input.r.Context().Value(CtxKeyLogger).([]error)
			assert.True(t, ok)
			assert.Equal(t, c.output.errs, errs)
		})
	}
}
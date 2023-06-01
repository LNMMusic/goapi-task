package task

import (
	"testing"

	"github.com/LNMMusic/optional"

	"github.com/stretchr/testify/assert"
)

// Tests
func TestValidatorLocal_Validate(t *testing.T) {
	type input struct {task *Task}
	type output struct {err error; errMsg string}
	type testCase struct {
		title  string
		input  input
		output output
	}

	cases := []testCase{
		// succeed cases
		{
			title: "valid task I",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(false),
			}},
			output: output{err: nil, errMsg: ""},
		},
		{
			title: "valid task II",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.None[string](),
				Completed: optional.Some(true),
			}},
			output: output{err: nil, errMsg: ""},
		},

		// failure cases
		{
			title: "invalid task - title required",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.None[string](),
				Description: optional.Some("description"),
				Completed: optional.Some(false),
			}},
			output: output{err: ErrValidatorFieldRequired, errMsg: "validator field required: title"},
		},
		{
			title: "invalid task - completed required",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.None[bool](),
			}},
			output: output{err: ErrValidatorFieldRequired, errMsg: "validator field required: completed"},
		},
		{
			title: "invalid task - title empty",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some(""),
				Description: optional.Some("description"),
				Completed: optional.Some(false),
			}},
			output: output{err: ErrValidatorFieldEmpty, errMsg: "validator field empty: title"},
		},
		{
			title: "invalid task - title quality",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title length is more than 50 characters so it is invalid"),
				Description: optional.Some("description"),
				Completed: optional.Some(false),
			}},
			output: output{err: ErrValidatorFieldQuality, errMsg: "validator field quality: title"},
		},
		{
			title: "invalid task - description quality",
			input: input{task: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description length is more than 150 characters so it is invalid and here are some more characters to make it invalid and even more characters to make it even more invalid"),
				Completed: optional.Some(false),
			}},
			output: output{err: ErrValidatorFieldQuality, errMsg: "validator field quality: description"},
		},
	}

	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			vl := NewValidatorLocal()

			// act
			err := vl.Validate(c.input.task)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if err != nil {
				assert.Equal(t, c.output.errMsg, err.Error())
			}
		})
	}
}
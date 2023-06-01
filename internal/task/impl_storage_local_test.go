package task

import (
	"fmt"
	"testing"

	"github.com/LNMMusic/optional"

	"github.com/stretchr/testify/assert"
)

// Tests
func TestStorageLocal_Get(t *testing.T) {
	type input struct {id string}
	type output struct {task *Task; err error; errMsg string}
	type testCase struct {
		title		 string
		input		 input
		output		 output
		setDatabase  func(db *[]*Task)
		setValidator func(vl *ValidatorMock)
	}

	cases := []testCase{
		// succeed cases
		{
			title: "get a task",
			input: input{id: "1"},
			output: output{
				task: &Task{
					ID: optional.Some("1"),
					Title: optional.Some("title"),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				},
				err: nil,
				errMsg: "",
			},
			setDatabase: func(db *[]*Task) {
				*db = []*Task{
					{
						ID: optional.Some("1"),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					},
				}
			},
			setValidator: func(vl *ValidatorMock) {},
		},

		// fail cases
		{
			title: "get a task that does not exist",
			input: input{id: "1"},
			output: output{
				task: nil,
				err: ErrStorageNotFound,
				errMsg: "storage task not found: 1",
			},
			setDatabase: func(db *[]*Task) {},
			setValidator: func(vl *ValidatorMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			db := []*Task{}
			c.setDatabase(&db)

			vl := NewValidatorMock()
			c.setValidator(vl)

			st := NewStorageLocal(db, vl)

			// act
			task, err := st.Get(c.input.id)

			// assert
			assert.Equal(t, c.output.task, task)
			assert.ErrorIs(t, err, c.output.err)
			if err != nil {
				assert.Equal(t, c.output.errMsg, err.Error())
			}
			vl.AssertExpectations(t)
		})
	}
}

func TestStorageLocal_Save(t *testing.T) {
	type input struct {task *Task}
	type output struct {err error; errMsg string}
	type testCase struct {
		title		 string
		input		 input
		output		 output
		setDatabase  func(db *[]*Task)
		setValidator func(vl *ValidatorMock)
	}

	cases := []testCase{
		// succeed cases
		{
			title: "save a task",
			input: input{
				task: &Task{
					ID: optional.None[string](),
					Title: optional.Some("title"),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				},
			},
			output: output{
				err: nil,
				errMsg: "",
			},
			setDatabase: func(db *[]*Task) {},
			setValidator: func(vl *ValidatorMock) {
				vl.On("Validate", &Task{
					ID: optional.None[string](),
					Title: optional.Some("title"),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				}).Return(nil)
			},
		},

		// failure cases
		{
			title: "save an invalid task",
			input: input{
				task: &Task{
					ID: optional.None[string](),
					Title: optional.None[string](),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				},
			},
			output: output{
				err: ErrStorageInvalid,
				errMsg: "storage invalid task: validation failed: title: is required",
			},
			setDatabase: func(db *[]*Task) {},
			setValidator: func(vl *ValidatorMock) {
				vl.On("Validate", &Task{
					ID: optional.None[string](),
					Title: optional.None[string](),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				}).Return(fmt.Errorf("validation failed: title: is required"))
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			db := []*Task{}
			c.setDatabase(&db)

			vl := NewValidatorMock()
			c.setValidator(vl)

			st := NewStorageLocal(db, vl)

			// act
			err := st.Save(c.input.task)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if err != nil {
				assert.Equal(t, c.output.errMsg, err.Error())
			}
			vl.AssertExpectations(t)
		})
	}
}
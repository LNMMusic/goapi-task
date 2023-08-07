package profiles

import (
	"testing"

	"github.com/LNMMusic/optional"
	"github.com/stretchr/testify/assert"
)

// Tests for ImplValidatorDefault
func TestImplValidatorDefault_Validate(t *testing.T) {
	type input struct { pf *Profile }
	type output struct { err error; errMsg string }
	type test struct {
		name string
		input input
		output output
	}

	cases := []test{
		// valid cases
		{
			name: "valid case - only required and valid fields",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
				},
			},
			output: output{err: nil, errMsg: ""},
		},
		{
			name: "valid case - all fields",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Name: optional.Some("name"),
					Email: optional.Some("johndoe@gmail.com"),
					Phone: optional.Some("1234567890"),
					Address: optional.Some("address"),
				},
			},
			output: output{err: nil, errMsg: ""},
		},

		// invalid cases
		// -> required fields
		{
			name: "invalid case - missing id",
			input: input{
				pf: &Profile{
					UserID: optional.Some("user_id"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - id field is required"},
		},
		{
			name: "invalid case - missing user_id",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - user_id field is required"},
		},
		// -> quality validation
		{
			name: "invalid case - user_id field empty",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some(""),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - user_id field can not be empty"},
		},
		{
			name: "invalid case - name field too short",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Name: optional.Some("na"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - name field must be between 3 and 50 characters"},
		},
		{
			name: "invalid case - name field too long",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Name: optional.Some("this name is way extremely too long and will cause an error"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - name field must be between 3 and 50 characters"},
		},
		{
			name: "invalid case - email field invalid",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Email: optional.Some("johndoegmail.com"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - email field is invalid"},
		},
		{
			name: "invalid case - phone field invalid",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Phone: optional.Some("123456789"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - phone field is invalid"},
		},
		{
			name: "invalid case - address field too short",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Address: optional.Some("ad"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - address field must be between 3 and 50 characters"},
		},
		{
			name: "invalid case - address field too long",
			input: input{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Address: optional.Some("this address is way extremely too long and will cause an error"),
				},
			},
			output: output{err: ErrValidatorInvalidProfile, errMsg: "validator: invalid profile - address field must be between 3 and 50 characters"},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			impl := NewImplValidatorDefault(nil)

			// act
			err := impl.Validate(c.input.pf)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
		})
	}
}
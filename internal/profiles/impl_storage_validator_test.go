package profiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for ImplStorageValidator
func TestImplStorageValidator_GetProfileByUserId(t *testing.T) {
	type input struct { userId string }
	type output struct { pf *Profile; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplStorageMock)
		setUpValidator func(mk *ImplValidatorMock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case",
			input: input{ userId: "user_id" },
			output: output{ pf: &Profile{}, err: nil, errMsg: "" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileByUserId", "user_id").Return(&Profile{}, nil)
			},
			setUpValidator: func(mk *ImplValidatorMock) {},
		},

		// invalid cases
		// -> storage
		{
			name: "storage error",
			input: input{ userId: "user_id" },
			output: output{ pf: &Profile{}, err: ErrStorageInternal, errMsg: "storage: internal storage error" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileByUserId", "user_id").Return(&Profile{}, ErrStorageInternal)
			},
			setUpValidator: func(mk *ImplValidatorMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplStorageMock()
			c.setUpStorage(st)

			vl := NewImplValidatorMock()
			c.setUpValidator(vl)

			impl := NewImplStorageValidator(st, vl)

			// act
			pf, err := impl.GetProfileByUserId(c.input.userId)

			// assert
			assert.Equal(t, c.output.pf, pf)
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			st.AssertExpectations(t)
			vl.AssertExpectations(t)
		})
	}
}

func TestImplStorageValidator_ActivateProfile(t *testing.T) {
	type input struct { pf *Profile }
	type output struct { err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplStorageMock)
		setUpValidator func(mk *ImplValidatorMock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case",
			input: input{ pf: &Profile{} },
			output: output{ err: nil, errMsg: "" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(nil)
			},
			setUpValidator: func(mk *ImplValidatorMock) {
				mk.On("Validate", &Profile{}).Return(nil)
			},
		},

		// invalid cases
		// -> validator
		{
			name: "validator error",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInvalidProfile, errMsg: "storage: invalid profile. validator: internal validator error" },
			setUpStorage: func(mk *ImplStorageMock) {},
			setUpValidator: func(mk *ImplValidatorMock) {
				mk.On("Validate", &Profile{}).Return(ErrValidatorInternal)
			},
		},
		// -> storage
		{
			name: "storage error",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInternal, errMsg: "storage: internal storage error" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(ErrStorageInternal)
			},
			setUpValidator: func(mk *ImplValidatorMock) {
				mk.On("Validate", &Profile{}).Return(nil)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplStorageMock()
			c.setUpStorage(st)

			vl := NewImplValidatorMock()
			c.setUpValidator(vl)

			impl := NewImplStorageValidator(st, vl)

			// act
			err := impl.ActivateProfile(c.input.pf)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			st.AssertExpectations(t)
			vl.AssertExpectations(t)
		})
	}
}
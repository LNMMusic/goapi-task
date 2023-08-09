package storage

import (
	"api/internal/profiles"
	"api/internal/profiles/validator"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for ImplProfilesStorageValidator
func TestImplProfilesStorageValidator_GetProfileById(t *testing.T) {
	type input struct { id string }
	type output struct { pf *profiles.Profile; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplProfilesStorageMock)
		setUpValidator func(mk *validator.ImplProfilesValidatorMock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case",
			input: input{ id: "id" },
			output: output{ pf: &profiles.Profile{}, err: nil, errMsg: "" },
			setUpStorage: func(mk *ImplProfilesStorageMock) {
				mk.On("GetProfileById", "id").Return(&profiles.Profile{}, nil)
			},
			setUpValidator: func(mk *validator.ImplProfilesValidatorMock) {},
		},

		// invalid cases
		// -> storage
		{
			name: "storage error",
			input: input{ id: "id" },
			output: output{ pf: &profiles.Profile{}, err: ErrStorageInternal, errMsg: "storage: internal storage error" },
			setUpStorage: func(mk *ImplProfilesStorageMock) {
				mk.On("GetProfileById", "id").Return(&profiles.Profile{}, ErrStorageInternal)
			},
			setUpValidator: func(mk *validator.ImplProfilesValidatorMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplProfilesStorageMock()
			c.setUpStorage(st)

			vl := validator.NewImplProfilesValidatorMock()
			c.setUpValidator(vl)

			impl := NewImplProfilesStorageValidator(st, vl)

			// act
			pf, err := impl.GetProfileById(c.input.id)

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

func TestImplProfilesStorageValidator_ActivateProfile(t *testing.T) {
	type input struct { pf *profiles.Profile }
	type output struct { err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplProfilesStorageMock)
		setUpValidator func(mk *validator.ImplProfilesValidatorMock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case",
			input: input{ pf: &profiles.Profile{} },
			output: output{ err: nil, errMsg: "" },
			setUpStorage: func(mk *ImplProfilesStorageMock) {
				mk.On("ActivateProfile", &profiles.Profile{}).Return(nil)
			},
			setUpValidator: func(mk *validator.ImplProfilesValidatorMock) {
				mk.On("Validate", &profiles.Profile{}).Return(nil)
			},
		},

		// invalid cases
		// -> validator
		{
			name: "validator error",
			input: input{ pf: &profiles.Profile{} },
			output: output{ err: ErrStorageInvalidProfile, errMsg: "storage: invalid profile. validator: internal validator error" },
			setUpStorage: func(mk *ImplProfilesStorageMock) {},
			setUpValidator: func(mk *validator.ImplProfilesValidatorMock) {
				mk.On("Validate", &profiles.Profile{}).Return(validator.ErrValidatorInternal)
			},
		},
		// -> storage
		{
			name: "storage error",
			input: input{ pf: &profiles.Profile{} },
			output: output{ err: ErrStorageInternal, errMsg: "storage: internal storage error" },
			setUpStorage: func(mk *ImplProfilesStorageMock) {
				mk.On("ActivateProfile", &profiles.Profile{}).Return(ErrStorageInternal)
			},
			setUpValidator: func(mk *validator.ImplProfilesValidatorMock) {
				mk.On("Validate", &profiles.Profile{}).Return(nil)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplProfilesStorageMock()
			c.setUpStorage(st)

			vl := validator.NewImplProfilesValidatorMock()
			c.setUpValidator(vl)

			impl := NewImplProfilesStorageValidator(st, vl)

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
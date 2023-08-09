package profiles

import (
	"api/pkg/mysql/transactioner"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests for ImplStorageMySQLTx
func TestImplStorageMySQLTx_GetProfileById(t *testing.T) {
	type input struct { id string }
	type output struct { pf *Profile; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplStorageMock)
		setUpTransactioner func(mk *transactioner.ImplTransactionerMock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case",
			input: input{ id: "id" },
			output: output{ pf: &Profile{}, err: nil, errMsg: "" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileById", "id").Return(&Profile{}, nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(nil)
			},
		},
		
		// invalid cases
		// -> operation error
		{
			name: "operation error - not found",
			input: input{ id: "id" },
			output: output{ pf: &Profile{}, err: ErrStorageNotFound, errMsg: "storage: profile not found" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileById", "id").Return(&Profile{}, ErrStorageNotFound)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionOperation)
			},
		},
		// -> default error
		{
			name: "default error - begin transaction",
			input: input{ id: "id" },
			output: output{ pf: &Profile{}, err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot begin transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileById", "id").Return(&Profile{}, nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionBegin)
			},
		},
		{
			name: "default error - commit transaction",
			input: input{ id: "id" },
			output: output{ pf: &Profile{}, err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot commit transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileById", "id").Return(&Profile{}, nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionCommit)
			},
		},
		{
			name: "default error - rollback transaction",
			input: input{ id: "id" },
			output: output{ pf: &Profile{}, err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot rollback transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("GetProfileById", "id").Return(&Profile{}, nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionRollback)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplStorageMock()
			c.setUpStorage(st)

			tr := transactioner.NewImplTransactionerMock()
			c.setUpTransactioner(tr)

			impl := NewImplStorageMySQLTx(st, tr)

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
			tr.AssertExpectations(t)
		})
	}
}

func TestImplStorageMySQLTx_ActivateProfile(t *testing.T) {
	type input struct { pf *Profile }
	type output struct { err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *ImplStorageMock)
		setUpTransactioner func(mk *transactioner.ImplTransactionerMock)
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
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(nil)
			},
		},

		// invalid cases
		// -> operation error
		{
			name: "operation error - invalid profile",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInvalidProfile, errMsg: "storage: invalid profile" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(ErrStorageInvalidProfile)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionOperation)
			},
		},
		// -> default error
		{
			name: "default error - begin transaction",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot begin transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionBegin)
			},
		},
		{
			name: "default error - commit transaction",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot commit transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionCommit)
			},
		},
		{
			name: "default error - rollback transaction",
			input: input{ pf: &Profile{} },
			output: output{ err: ErrStorageInternal, errMsg: "storage: internal storage error. transactioner: cannot rollback transaction" },
			setUpStorage: func(mk *ImplStorageMock) {
				mk.On("ActivateProfile", &Profile{}).Return(nil)
			},
			setUpTransactioner: func(mk *transactioner.ImplTransactionerMock) {
				mk.On("Do", mock.Anything).Return(transactioner.ErrTransactionRollback)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := NewImplStorageMock()
			c.setUpStorage(st)

			tr := transactioner.NewImplTransactionerMock()
			c.setUpTransactioner(tr)

			impl := NewImplStorageMySQLTx(st, tr)

			// act
			err := impl.ActivateProfile(c.input.pf)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			st.AssertExpectations(t)
			tr.AssertExpectations(t)
		})
	}
}
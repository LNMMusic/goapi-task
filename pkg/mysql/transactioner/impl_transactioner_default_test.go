package transactioner

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Tests for Do function
func TestImplTransactionerDefault(t *testing.T) {
	type input struct { op operation }
	type output struct { err error; errMsg string }
	type test struct {
		// base
		name string
		input input
		output output
		// process
		setUpMockDB func(mk sqlmock.Sqlmock)
	}

	cases := []test{
		// valid case
		{
			name: "valid case",
			input: input{op: func() (err error) {return}},
			output: output{err: nil, errMsg: ""},
			setUpMockDB: func(mk sqlmock.Sqlmock) {
				mk.ExpectBegin()
				mk.ExpectCommit()
			},
		},

		// invalid case
		// -> begin transaction error
		{
			name: "begin transaction error",
			input: input{op: func() (err error) {return}},
			output: output{err: ErrTransactionBegin, errMsg: "transactioner: cannot begin transaction. mysql begin error"},
			setUpMockDB: func(mk sqlmock.Sqlmock) {
				mk.ExpectBegin().WillReturnError(errors.New("mysql begin error"))
			},
		},
		// -> operation error
		{
			name: "operation error",
			input: input{op: func() (err error) {return errors.New("operation error")}},
			output: output{err: ErrTransactionOperation, errMsg: "transactioner: operation failed. operation error"},
			setUpMockDB: func(mk sqlmock.Sqlmock) {
				mk.ExpectBegin()
				mk.ExpectRollback()
			},
		},
		// -> rollback transaction error
		{
			name: "rollback transaction error",
			input: input{op: func() (err error) {return errors.New("operation error")}},
			output: output{err: ErrTransactionRollback, errMsg: "transactioner: cannot rollback transaction. mysql rollback error"},
			setUpMockDB: func(mk sqlmock.Sqlmock) {
				mk.ExpectBegin()
				mk.ExpectRollback().WillReturnError(errors.New("mysql rollback error"))
			},
		},
		// -> commit transaction error
		{
			name: "commit transaction error",
			input: input{op: func() (err error) {return}},
			output: output{err: ErrTransactionCommit, errMsg: "transactioner: cannot commit transaction. mysql commit error"},
			setUpMockDB: func(mk sqlmock.Sqlmock) {
				mk.ExpectBegin()
				mk.ExpectCommit().WillReturnError(errors.New("mysql commit error"))
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			db, mk, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			c.setUpMockDB(mk)

			impl := NewImplTransactionerDefault(db)

			// act
			err = impl.Do(c.input.op)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
		})
	}
}
package profiles

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LNMMusic/optional"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// Tests for ImplStorageMySQL
func TestImplStorageMySQL_GetProfileByUserId(t *testing.T) {
	type input struct { userId string }
	type output struct { pf *Profile; err error; errMsg string }
	type test struct {
		name string
		input input
		output output
		// set-up
		setUpDB func (mk sqlmock.Sqlmock)
	}

	cases := []test{
		// valid cases
		{
			name: "valid case - found",
			input: input{userId: "user_id"},
			output: output{
				pf: &Profile{
					ID: optional.Some("id"),
					UserID: optional.Some("user_id"),
					Name: optional.Some("name"),
					Email: optional.Some("johndoe@gmail.com"),
					Phone: optional.Some("1234567890"),
					Address: optional.Some("address"),
				},
				err: nil, errMsg: "",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE user_id = ?"
				
				cols := []string{"id", "user_id", "name", "email", "phone", "address"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "id", Valid: true},
					sql.NullString{String: "user_id", Valid: true},
					sql.NullString{String: "name", Valid: true},
					sql.NullString{String: "johndoe@gmail.com", Valid: true},
					sql.NullString{String: "1234567890", Valid: true},
					sql.NullString{String: "address", Valid: true},
				)

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().WithArgs("user_id").
					WillReturnRows(rows)
			},
		},

		// invalid cases
		// -> query error. no rows
		{
			name: "invalid case - not found",
			input: input{userId: "user_id"},
			output: output{
				pf: nil,
				err: ErrStorageNotFound, errMsg: "storage: profile not found. sql: no rows in result set",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE user_id = ?"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().WithArgs("user_id").
					WillReturnError(sql.ErrNoRows)
			},
		},
		// -> query error. internal error
		{
			name: "invalid case - scan internal error",
			input: input{userId: "user_id"},
			output: output{
				pf: nil,
				err: ErrStorageInternal, errMsg: "storage: internal storage error. sql: internal error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE user_id = ?"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectQuery().WithArgs("user_id").
					WillReturnError(errors.New("sql: internal error"))
			},
		},
		// -> prepare error
		{
			name: "invalid case - prepare internal error",
			input: input{userId: "user_id"},
			output: output{
				pf: nil,
				err: ErrStorageInternal, errMsg: "storage: internal storage error. sql: prepare error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE user_id = ?"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("sql: prepare error"))
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

			c.setUpDB(mk)

			impl := NewImplStorageMySQL(db)

			// act
			pf, err := impl.GetProfileByUserId(c.input.userId)

			// assert
			assert.Equal(t, c.output.pf, pf)
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
		})
	}
}

func TestImplStorageMySQL_ActiveProfile(t *testing.T) {
	type input struct { pf *Profile }
	type output struct { err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpDB func (mk sqlmock.Sqlmock)
	}

	cases := []testCase{
		// valid cases
		{
			name: "valid case - success",
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
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WithArgs(
						sql.NullString{String: "id", Valid: true},
						sql.NullString{String: "user_id", Valid: true},
						sql.NullString{String: "name", Valid: true},
						sql.NullString{String: "johndoe@gmail.com", Valid: true},
						sql.NullString{String: "1234567890", Valid: true},
						sql.NullString{String: "address", Valid: true},
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},

		// invalid cases
		// -> prepare error
		{
			name: "invalid case - prepare internal error",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageInternal, errMsg: "storage: internal storage error. sql: prepare error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					WillReturnError(errors.New("sql: prepare error"))
			},
		},
		// -> exec error. default error
		{
			name: "invalid case - exec internal error",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageInternal, errMsg: "storage: internal storage error. sql: exec error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().
					WillReturnError(errors.New("sql: exec error"))
			},
		},
		// -> exec error. mysql error - duplicate entry
		{
			name: "invalid case - exec duplicate entry error",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageNotUnique, errMsg: "storage: profile not unique. Error 1062: duplicate entry",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().
					WillReturnError(&mysql.MySQLError{Number: 1062, Message: "duplicate entry"})
			},
		},
		// -> exec error. mysql error - other
		{
			name: "invalid case - exec mysql error",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageInternal, errMsg: "storage: internal storage error. Error 1234: other error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().
					WillReturnError(&mysql.MySQLError{Number: 1234, Message: "other error"})
			},
		},
		// -> result error.
		{
			name: "invalid case - result error",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageInternal, errMsg: "storage: internal storage error. sql: result error",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().
					WillReturnResult(sqlmock.NewErrorResult(errors.New("sql: result error")))
			},
		},
		// -> result error. rows affected != 1
		{
			name: "invalid case - result error. rows affected != 1",
			input: input{pf: &Profile{}},
			output: output{
				err: ErrStorageInternal, errMsg: "storage: internal storage error. rows affected != 1",
			},
			setUpDB: func (mk sqlmock.Sqlmock) {
				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				// expectations
				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 0))
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

			c.setUpDB(mk)

			impl := NewImplStorageMySQL(db)

			// act
			err = impl.ActivateProfile(c.input.pf)

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
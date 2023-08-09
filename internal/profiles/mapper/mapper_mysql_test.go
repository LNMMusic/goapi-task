package mapper

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Tests for ProfileMapperMySQL implementation
func TestProfileMapperMySQL_MapProfile(t *testing.T) {
	type input struct { userId string }
	type output struct { profileId string; err error; errMsg string }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpDatabase func (mk sqlmock.Sqlmock)
	}

	cases := []testCase{
		// valid cases
		// -> profile found
		{
			name: "valid case - profile found",
			input: input{ userId: "user-id-1" },
			output: output{ profileId: "profile-id-1", err: nil, errMsg: "" },
			setUpDatabase: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT profile_id FROM profiles WHERE user_id = ?"
				
				cols := []string{"profile_id"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow("profile-id-1")

				// statement
				mk.
					ExpectPrepare(query).
					ExpectQuery().WithArgs("user-id-1").
					WillReturnRows(rows)
			},
		},

		// error cases
		// -> prepare error
		{
			name: "error case - prepare error",
			input: input{ userId: "user-id-1" },
			output: output{ profileId: "", err: ErrProfileMapperInternal, errMsg: "mapper: internal mapper error. prepare error" },
			setUpDatabase: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT profile_id FROM profiles WHERE user_id = ?"

				// statement
				mk.
					ExpectPrepare(query).
					WillReturnError(errors.New("prepare error"))
			},
		},
		// -> query error - no rows
		{
			name: "error case - query error - no rows",
			input: input{ userId: "user-id-1" },
			output: output{ profileId: "", err: ErrProfileMapperNotFound, errMsg: "mapper: mapper not found. sql: no rows in result set" },
			setUpDatabase: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT profile_id FROM profiles WHERE user_id = ?"

				// statement
				mk.
					ExpectPrepare(query).
					ExpectQuery().WithArgs("user-id-1").
					WillReturnError(sql.ErrNoRows)
			},
		},
		// -> query error - default
		{
			name: "error case - query error - default",
			input: input{ userId: "user-id-1" },
			output: output{ profileId: "", err: ErrProfileMapperInternal, errMsg: "mapper: internal mapper error. query error default" },
			setUpDatabase: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT profile_id FROM profiles WHERE user_id = ?"

				// statement
				mk.
					ExpectPrepare(query).
					ExpectQuery().WithArgs("user-id-1").
					WillReturnError(errors.New("query error default"))
			},
		},
		// -> scan error
		{
			name: "error case - scan error",
			input: input{ userId: "user-id-1" },
			output: output{ profileId: "", err: ErrProfileMapperInternal, errMsg: "mapper: internal mapper error. sql: Scan error on column index 0, name \"profile_id\": converting NULL to string is unsupported" },
			setUpDatabase: func (mk sqlmock.Sqlmock) {
				// query
				query := "SELECT profile_id FROM profiles WHERE user_id = ?"

				cols := []string{"profile_id"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(nil)
				
				// statement
				mk.
					ExpectPrepare(query).
					ExpectQuery().WithArgs("user-id-1").
					WillReturnRows(rows)
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
			
			c.setUpDatabase(mk)

			impl := NewProfileMapperMySQL(db)

			// act
			profileId, err := impl.MapProfile(c.input.userId)

			// assert
			assert.Equal(t, c.output.profileId, profileId)
			assert.ErrorIs(t, err, c.output.err)
			if c.output.err != nil {
				assert.EqualError(t, err, c.output.errMsg)
			}
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
		})
	}
}
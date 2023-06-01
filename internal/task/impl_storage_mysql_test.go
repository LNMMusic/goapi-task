package task

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LNMMusic/optional"
	"github.com/stretchr/testify/assert"
)

// Tests
func TestStorageMySQL_Get(t *testing.T) {
	type input struct {id string}
	type output struct {ts *Task; err error; errMsg string}
	type testCase struct {
		// io
		title  		 string
		input  		 input
		output 		 output
		// process
		setDatabase  func(mk sqlmock.Sqlmock)
		setValidator func(mk *ValidatorMock)
	}

	cases := []testCase{
		// success cases
		{
			title: "full task",
			input: input{id: "id"},
			output: output{
				ts: &Task{
					ID: optional.Some("id"),
					Title: optional.Some("title"),
					Description: optional.Some("description"),
					Completed: optional.Some(true),
				},
				err: nil,
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// rows
				cols := []string{"id", "title", "description", "completed"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "id", Valid: true},
					sql.NullString{String: "title", Valid: true},
					sql.NullString{String: "description", Valid: true},
					sql.NullBool{Bool: true, Valid: true},
				)

				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					ExpectQuery().WithArgs("id").
					WillReturnRows(rows)
			},
			setValidator: func(mk *ValidatorMock) {},
		},
		{
			title: "empty task",
			input: input{id: "id"},
			output: output{
				ts: &Task{
					ID: optional.Some("id"),
					Title: optional.None[string](),
					Description: optional.None[string](),
					Completed: optional.None[bool](),
				},
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// rows
				cols := []string{"id", "title", "description", "completed"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "id", Valid: true},
					sql.NullString{String: "", Valid: false},
					sql.NullString{String: "", Valid: false},
					sql.NullBool{Bool: false, Valid: false},
				)

				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					ExpectQuery().WithArgs("id").
					WillReturnRows(rows)
			},
			setValidator: func(mk *ValidatorMock) {},
		},
		{
			title: "task with some fields",
			input: input{id: "id"},
			output: output{
				ts: &Task{
					ID: optional.Some("id"),
					Title: optional.Some("title"),
					Description: optional.None[string](),
					Completed: optional.None[bool](),
				},
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// rows
				cols := []string{"id", "title", "description", "completed"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "id", Valid: true},
					sql.NullString{String: "title", Valid: true},
					sql.NullString{String: "", Valid: false},
					sql.NullBool{Bool: false, Valid: false},
				)

				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					ExpectQuery().WithArgs("id").
					WillReturnRows(rows)
			},
			setValidator: func(mk *ValidatorMock) {},
		},

		// failure cases
		{
			title: "prepare error",
			input: input{id: "id"},
			output: output{
				ts: nil,
				err: ErrStorageInternal,
				errMsg: "storage internal error: prepare",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					WillReturnError(sql.ErrConnDone)
			},
			setValidator: func(mk *ValidatorMock) {},
		},
		{
			title: "non existing task",
			input: input{id: "id"},
			output: output{
				ts: nil,
				err: ErrStorageNotFound,
				errMsg: "storage task not found: query row",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					ExpectQuery().WithArgs("id").
					WillReturnError(sql.ErrNoRows)
			},
			setValidator: func(mk *ValidatorMock) {},
		},
		{
			title: "database error",
			input: input{id: "id"},
			output: output{
				ts: nil,
				err: ErrStorageInternal,
				errMsg: "storage internal error: query row",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				mk.
					ExpectPrepare(regexp.QuoteMeta(QueryGetTask)).
					ExpectQuery().WithArgs("id").
					WillReturnError(sql.ErrConnDone)
			},
			setValidator: func(mk *ValidatorMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			db, mk, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			c.setDatabase(mk)

			vl := NewValidatorMock()
			c.setValidator(vl)

			st := NewStorageMySQL(db, vl)

			// act
			ts, err := st.Get(c.input.id)

			// assert
			assert.Equal(t, c.output.ts, ts)
			assert.ErrorIs(t, err, c.output.err)
			if err != nil {
				assert.Equal(t, c.output.errMsg, err.Error())
			}
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
			vl.AssertExpectations(t)
		})
	}
}

func TestStorageMySQL_Save(t *testing.T) {
	type input struct {ts *Task}
	type output struct {err error; errMsg string}
	type testCase struct {
		// io
		title  		 string
		input  		 input
		output 		 output
		// process
		setDatabase  func(mk sqlmock.Sqlmock)
		setValidator func(mk *ValidatorMock)
	}

	cases := []testCase{
		// success cases
		{
			title: "full task",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{err: nil, errMsg: ""},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				// -> begin
				mk.ExpectBegin()

				// -> stmt
				mk.
					ExpectPrepare(regexp.QuoteMeta(QuerySaveTask)).
					ExpectExec().WithArgs(
						sqlmock.AnyArg(),
						sql.NullString{String: "title", Valid: true},
						sql.NullString{String: "description", Valid: true},
						sql.NullBool{Bool: true, Valid: true},							
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				
				// -> commit
				mk.ExpectCommit()
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},

		// failure cases
		// -> validator
		{
			title: "validator error",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.None[string](),
				Description: optional.None[string](),
				Completed: optional.None[bool](),
			}},
			output: output{
				err: ErrStorageInvalid,
				errMsg: "storage invalid task: validate",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.None[string](),
						Description: optional.None[string](),
						Completed: optional.None[bool](),
					}).
					Return(ErrStorageInvalid)
			},
		},
		// -> database
		{
			title: "init transaction error",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{
				err: ErrStorageInternal,
				errMsg: "storage internal error: begin",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				mk.ExpectBegin().WillReturnError(sql.ErrConnDone)
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},
		{
			title: "prepare statement error",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{
				err: ErrStorageInternal,
				errMsg: "storage internal error: prepare",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				// -> begin
				mk.ExpectBegin()

				// -> stmt
				mk.
					ExpectPrepare(regexp.QuoteMeta(QuerySaveTask)).
					WillReturnError(sql.ErrConnDone)

				// -> commit
				mk.ExpectRollback()
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},
		{
			title: "execute statement error",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{
				err: ErrStorageInternal,
				errMsg: "storage internal error: exec",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				// -> begin
				mk.ExpectBegin()

				// -> stmt
				mk.
					ExpectPrepare(regexp.QuoteMeta(QuerySaveTask)).
					ExpectExec().WithArgs(
						sqlmock.AnyArg(),
						sql.NullString{String: "title", Valid: true},
						sql.NullString{String: "description", Valid: true},
						sql.NullBool{Bool: true, Valid: true},							
					).
					WillReturnError(sql.ErrConnDone)

				// -> commit
				mk.ExpectRollback()
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},
		{
			title: "rows affected error",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{
				err: ErrStorageInternal,
				errMsg: "storage internal error: result rows affected",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				// -> begin
				mk.ExpectBegin()

				// -> stmt
				mk.
					ExpectPrepare(regexp.QuoteMeta(QuerySaveTask)).
					ExpectExec().WithArgs(
						sqlmock.AnyArg(),
						sql.NullString{String: "title", Valid: true},
						sql.NullString{String: "description", Valid: true},
						sql.NullBool{Bool: true, Valid: true},							
					).
					WillReturnResult(sqlmock.NewErrorResult(sql.ErrConnDone))

				// -> commit
				mk.ExpectRollback()
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},
		{
			title: "rows affected not 1",
			input: input{ts: &Task{
				ID: optional.None[string](),
				Title: optional.Some("title"),
				Description: optional.Some("description"),
				Completed: optional.Some(true),
			}},
			output: output{
				err: ErrStorageInternal,
				errMsg: "storage internal error: rows affected",
			},
			setDatabase: func(mk sqlmock.Sqlmock) {
				// mock
				// -> begin
				mk.ExpectBegin()

				// -> stmt
				mk.
					ExpectPrepare(regexp.QuoteMeta(QuerySaveTask)).
					ExpectExec().WithArgs(
						sqlmock.AnyArg(),
						sql.NullString{String: "title", Valid: true},
						sql.NullString{String: "description", Valid: true},
						sql.NullBool{Bool: true, Valid: true},							
					).
					WillReturnResult(sqlmock.NewResult(1, 0))

				// -> commit
				mk.ExpectRollback()
			},
			setValidator: func(mk *ValidatorMock) {
				mk.
					On("Validate", &Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(true),
					}).
					Return(nil)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			db, mk, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()
			c.setDatabase(mk)

			vl := NewValidatorMock()
			c.setValidator(vl)

			st := NewStorageMySQL(db, vl)

			// act
			err = st.Save(c.input.ts)

			// assert
			assert.ErrorIs(t, err, c.output.err)
			if err != nil {
				assert.Equal(t, c.output.errMsg, err.Error())
			}
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
			vl.AssertExpectations(t)
		})
	}
}
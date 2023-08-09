package handlers

import (
	"api/internal/profiles/contexter"
	"api/internal/profiles/storage"
	"api/internal/profiles/validator"
	"api/pkg/mysql/transactioner"
	"api/pkg/uuidgenerator"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

// Functional Tests for Profile Controller
func TestFunctionalProfileController_GetProfileByUserId(t *testing.T) {
	type input struct { r *http.Request; rr *httptest.ResponseRecorder; setR func (r *http.Request) }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		// -> mocks
		setUpDatabase func (mk sqlmock.Sqlmock)
		setUpUUID func (mk *uuidgenerator.ImplUUIDGeneratorMock)
		// -> impl
		// setUpValidator func (cfg *profiles.Config) // default config
	}

	cases := []testCase{
		// success to get profile by user id
		{
			name: "success to get profile by user id",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
				setR: func (r *http.Request) {
					// set-up context
					(*r) = *(*r).WithContext(context.WithValue((*r).Context(), contexter.KeyProfileId, "1"))
				},
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Success","data":{"user_id":"1","name":"John Doe","email":"johndoe@gmail.com", "phone":"111122223", "address":"Jl. Raya Bogor"}, "error":false}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE id = ?" 

				cols := []string{"id", "user_id", "name", "email", "phone", "address"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "1", Valid: true},
					sql.NullString{String: "1", Valid: true},
					sql.NullString{String: "John Doe", Valid: true},
					sql.NullString{String: "johndoe@gmail.com", Valid: true},
					sql.NullString{String: "111122223", Valid: true},
					sql.NullString{String: "Jl. Raya Bogor", Valid: true},
				)

				mk.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("1").WillReturnRows(rows)

				// commit
				mk.ExpectCommit()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {},
		},
		// success to get profile by user id - missing fields
		{
			name: "success to get profile by user id - missing fields",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
				setR: func (r *http.Request) {
					// set-up context
					(*r) = *(*r).WithContext(context.WithValue((*r).Context(), contexter.KeyProfileId, "1"))
				},
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Success","data":{"user_id":"1","name":null,"email":null, "phone":null, "address":null}, "error":false}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE id = ?" 

				cols := []string{"id", "user_id", "name", "email", "phone", "address"}
				rows := sqlmock.NewRows(cols)
				rows.AddRow(
					sql.NullString{String: "1", Valid: true},
					sql.NullString{String: "1", Valid: true},
					sql.NullString{String: "", Valid: false},
					sql.NullString{String: "", Valid: false},
					sql.NullString{String: "", Valid: false},
					sql.NullString{String: "", Valid: false},
				)

				mk.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("1").WillReturnRows(rows)

				// commit
				mk.ExpectCommit()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {},
		},

		// fail to get profile by user id - profile not found
		{
			name: "fail to get profile by user id - profile not found",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
				setR: func (r *http.Request) {
					// set-up context
					(*r) = *(*r).WithContext(context.WithValue((*r).Context(), contexter.KeyProfileId, "1"))
				},
			},
			output: output{
				code: http.StatusNotFound,
				body: `{"message":"Profile not found","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE id = ?" 

				mk.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("1").WillReturnError(sql.ErrNoRows)

				// rollback
				mk.ExpectRollback()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {},
		},
		// fail to get profile by user id - internal server error
		{
			name: "fail to get profile by user id - internal server error",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
				setR: func (r *http.Request) {
					// set-up context
					(*r) = *(*r).WithContext(context.WithValue((*r).Context(), contexter.KeyProfileId, "1"))
				},
			},
			output: output{
				code: http.StatusInternalServerError,
				body: `{"message":"Internal server error","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE id = ?" 

				mk.ExpectPrepare(regexp.QuoteMeta(query)).ExpectQuery().WithArgs("1").WillReturnError(sql.ErrConnDone)

				// rollback
				mk.ExpectRollback()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {},
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

			vl := validator.NewImplProfilesValidatorDefault(&validator.Config{})
			tx := transactioner.NewImplTransactionerDefault(db)
			st := storage.NewImplProfilesStorageValidator(
				storage.NewImplProfilesStorageMySQLTx(
					storage.NewImplProfilesStorageMySQL(db),
					tx,
				),
				vl,
			)
			uuid := uuidgenerator.NewUUIDGeneratorMock()
			c.setUpUUID(uuid)

			ct := NewProfileController(st, uuid)
			hd := ct.GetProfileById()

			// act
			c.input.setR(c.input.r)
			hd(c.input.rr, c.input.r)

			// assert
			assert.Equal(t, c.output.code, c.input.rr.Code)
			assert.JSONEq(t, c.output.body, c.input.rr.Body.String())
			assert.Equal(t, c.output.headers, c.input.rr.Header())
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
			uuid.AssertExpectations(t)
		})
	}
}

func TestFunctionalProfileController_ActivateProfile(t *testing.T) {
	type input struct { r *http.Request; rr *httptest.ResponseRecorder }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpDatabase func(mk sqlmock.Sqlmock)
		setUpUUID func(uuid *uuidgenerator.ImplUUIDGeneratorMock)
	}

	cases := []testCase{
		// success to activate profile
		{
			name: "success",
			input: input{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Success","data":null,"error":false}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WithArgs(
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
					).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// commit
				mk.ExpectCommit()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {
				uuid.On("UUID").Return("1")
			},
		},
		// fail to activate profile - invalid profile
		{
			name: "fail to activate profile - invalid profile",
			input: input{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": []string{""},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusUnprocessableEntity,
				body: `{"message":"Invalid profile","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {
				uuid.On("UUID").Return("1")
			},
		},
		// fail to activate profile - not unique user id
		{
			name: "fail to activate profile - not unique user id",
			input: input{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusConflict,
				body: `{"message":"Profile not unique","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WithArgs(
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
					).
					WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry '1' for key 'profiles.user_id_UNIQUE'"})

				// rollback
				mk.ExpectRollback()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {
				uuid.On("UUID").Return("1")
			},
		},
		// fail to activate profile - internal error
		{
			name: "fail to activate profile - internal error",
			input: input{
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": []string{"1"},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusInternalServerError,
				body: `{"message":"Internal server error","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpDatabase: func(mk sqlmock.Sqlmock) {
				// transaction
				mk.ExpectBegin()

				// query
				query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

				mk.
					ExpectPrepare(regexp.QuoteMeta(query)).
					ExpectExec().WithArgs(
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "1", Valid: true},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
						sql.NullString{String: "", Valid: false},
					).
					WillReturnError(errors.New("unexpected error"))

				// rollback
				mk.ExpectRollback()
			},
			setUpUUID: func(uuid *uuidgenerator.ImplUUIDGeneratorMock) {
				uuid.On("UUID").Return("1")
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

			vl := validator.NewImplProfilesValidatorDefault(&validator.Config{})
			tx := transactioner.NewImplTransactionerDefault(db)
			st := storage.NewImplProfilesStorageValidator(
				storage.NewImplProfilesStorageMySQLTx(
					storage.NewImplProfilesStorageMySQL(db),
					tx,
				),
				vl,
			)
			uuid := uuidgenerator.NewUUIDGeneratorMock()
			c.setUpUUID(uuid)

			ct := NewProfileController(st, uuid)
			hd := ct.ActivateProfile()

			// act
			hd(c.input.rr, c.input.r)
			
			// assert
			assert.Equal(t, c.output.code, c.input.rr.Code)
			assert.JSONEq(t, c.output.body, c.input.rr.Body.String())
			assert.Equal(t, c.output.headers, c.input.rr.Header())
			// -> expectations
			assert.NoError(t, mk.ExpectationsWereMet())
			uuid.AssertExpectations(t)
		})
	}
}
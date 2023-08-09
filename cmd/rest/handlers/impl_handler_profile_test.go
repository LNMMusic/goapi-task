package handlers

import (
	"api/internal/profiles"
	"api/internal/profiles/contexter"
	"api/internal/profiles/storage"
	"api/pkg/uuidgenerator"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LNMMusic/optional"
	"github.com/stretchr/testify/assert"
)

// Tests for ProfileController handlers
func TestProfileController_GetProfileById(t *testing.T) {
	type input struct { w *httptest.ResponseRecorder; r *http.Request; setR func (r *http.Request) }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *storage.ImplProfilesStorageMock)
		setUpUUID func(mk *uuidgenerator.ImplUUIDGeneratorMock)
	}

	cases := []testCase{
		// valid case
		{
			name: "valid case",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
				setR: func (r *http.Request) {
					// set-up request context
					(*r) = *(*r).WithContext(context.WithValue(r.Context(), contexter.KeyProfileId, "id"))
				},
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Success","data":{"user_id":"user_id","name":"name","email":"email","phone":"phone","address":"address"},"error":false}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("GetProfileById", "id").
					Return(&profiles.Profile{
						ID:      optional.Some("id"),
						UserID:  optional.Some("user_id"),
						Name:    optional.Some("name"),
						Email:   optional.Some("email"),
						Phone:   optional.Some("phone"),
						Address: optional.Some("address"),
					}, nil)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {},
		},

		// invalid case: storage error - not found
		{
			name: "invalid case: storage error",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
				setR: func (r *http.Request) {
					// set-up request context
					(*r) = *r.WithContext(context.WithValue(r.Context(), contexter.KeyProfileId, "id"))
				},
			},
			output: output{
				code: http.StatusNotFound,
				body: `{"message":"Profile not found","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("GetProfileById", "id").
					Return(&profiles.Profile{}, storage.ErrStorageNotFound)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {},
		},
		// invalid case: storage error - internal
		{
			name: "invalid case: storage error - internal",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
				setR: func (r *http.Request) {
					// set-up request context
					(*r) = *r.WithContext(context.WithValue(r.Context(), contexter.KeyProfileId, "id"))
				},
			},
			output: output{
				code: http.StatusInternalServerError,
				body: `{"message":"Internal server error","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("GetProfileById", "id").
					Return(&profiles.Profile{}, storage.ErrStorageInternal)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := storage.NewImplProfilesStorageMock()
			c.setUpStorage(st)

			uuid := uuidgenerator.NewUUIDGeneratorMock()
			c.setUpUUID(uuid)

			ct := NewProfileController(st, uuid)
			hd := ct.GetProfileById()

			// act
			c.input.setR(c.input.r)
			hd(c.input.w, c.input.r)

			// assert
			assert.Equal(t, c.output.code, c.input.w.Code)
			assert.JSONEq(t, c.output.body, c.input.w.Body.String())
			assert.Equal(t, c.output.headers, c.input.w.Header())
			// -> expectations
			st.AssertExpectations(t)
			uuid.AssertExpectations(t)
		})
	}
}

func TestProfileController_ActivateProfile(t *testing.T) {
	type input struct { w *httptest.ResponseRecorder; r *http.Request }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpStorage func(mk *storage.ImplProfilesStorageMock)
		setUpUUID func(mk *uuidgenerator.ImplUUIDGeneratorMock)
	}

	cases := []testCase{
		// valid case
		{
			name: "valid case",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Success","data":null,"error":false}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("ActivateProfile", &profiles.Profile{
						ID:      optional.Some("id"),
						UserID:  optional.Some("user_id"),
					}).
					Return(nil)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {
				mk.
					On("UUID").
					Return("id")
			},
		},

		// invalid case: storage error - not unique
		{
			name: "invalid case: storage error - not unique",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
			},
			output: output{
				code: http.StatusConflict,
				body: `{"message":"Profile not unique","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("ActivateProfile", &profiles.Profile{
						ID:      optional.Some("id"),
						UserID:  optional.Some("user_id"),
					}).
					Return(storage.ErrStorageNotUnique)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {
				mk.
					On("UUID").
					Return("id")
			},
		},

		// invalid case: storage error - internal
		{
			name: "invalid case: storage error - internal",
			input: input{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodPost,
					Header: http.Header{
						"User-Id": {"user_id"},
					},
				},
			},
			output: output{
				code: http.StatusInternalServerError,
				body: `{"message":"Internal server error","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": {"application/json"},
				},
			},
			setUpStorage: func(mk *storage.ImplProfilesStorageMock) {
				mk.
					On("ActivateProfile", &profiles.Profile{
						ID:      optional.Some("id"),
						UserID:  optional.Some("user_id"),
					}).
					Return(storage.ErrStorageInternal)
			},
			setUpUUID: func(mk *uuidgenerator.ImplUUIDGeneratorMock) {
				mk.
					On("UUID").
					Return("id")
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			st := storage.NewImplProfilesStorageMock()
			c.setUpStorage(st)

			uuid := uuidgenerator.NewUUIDGeneratorMock()
			c.setUpUUID(uuid)

			ct := NewProfileController(st, uuid)
			hd := ct.ActivateProfile()

			// act
			hd(c.input.w, c.input.r)

			// assert
			assert.Equal(t, c.output.code, c.input.w.Code)
			assert.JSONEq(t, c.output.body, c.input.w.Body.String())
			assert.Equal(t, c.output.headers, c.input.w.Header())
			// -> expectations
			st.AssertExpectations(t)
			uuid.AssertExpectations(t)
		})
	}
}
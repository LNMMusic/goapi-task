package mapping

import (
	"api/internal/profiles/mapper"
	"api/pkg/httpmock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests for ProfileMapping
func TestProfileMapping_MapProfile(t *testing.T) {
	type input struct { r *http.Request; rr *httptest.ResponseRecorder }
	type output struct { code int; body string; headers http.Header }
	type testCase struct {
		name string
		input input
		output output
		// set-up
		setUpMapperMock func(mk *mapper.ProfileMapperMock)
		setUpHandlerMock func(mk *httpmock.HandlerMock)
	}

	cases := []testCase{
		// valid case
		{
			name: "valid case - profile mapped",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"user-id"},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusOK,
				body: `{"message":"Profile mapped","data":null,"error":false}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpMapperMock: func(mk *mapper.ProfileMapperMock) {
				mk.On("MapProfile", "user-id").Return("profile-id", nil)
			},
			setUpHandlerMock: func(mk *httpmock.HandlerMock) {
				// set-up serveHTTP
				(*mk).SetUpServeHTTP = func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"message":"Profile mapped","data":null,"error":false}`))
					w.Header().Set("Content-Type", "application/json")
				}

				mk.On("ServeHTTP", mock.Anything, mock.Anything).Return()
			},
		},
		
		// invalid case
		// -> mapper error
		{
			name: "invalid case - mapper error - profile not found",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"user-id"},
					},
				},
				rr: httptest.NewRecorder(),
			},
			output: output{
				code: http.StatusUnauthorized,
				body: `{"message":"Profile not found","data":null,"error":true}`,
				headers: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			setUpMapperMock: func(mk *mapper.ProfileMapperMock) {
				mk.On("MapProfile", "user-id").Return("", mapper.ErrProfileMapperNotFound)
			},
			setUpHandlerMock: func(mk *httpmock.HandlerMock) {},
		},
		{
			name: "invalid case - mapper error - internal error",
			input: input{
				r: &http.Request{
					Method: http.MethodGet,
					Header: http.Header{
						"User-Id": []string{"user-id"},
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
			setUpMapperMock: func(mk *mapper.ProfileMapperMock) {
				mk.On("MapProfile", "user-id").Return("", mapper.ErrProfileMapperInternal)
			},
			setUpHandlerMock: func(mk *httpmock.HandlerMock) {},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// arrange
			mpMock := mapper.NewProfileMapperMock()
			c.setUpMapperMock(mpMock)

			impl := NewProfileMapping(mpMock)

			hdMock := httpmock.NewHandlerMock()
			c.setUpHandlerMock(hdMock)

			md := impl.MapProfile(hdMock)
			hd := md.(http.HandlerFunc)

			// act
			hd(c.input.rr, c.input.r)

			// assert
			assert.Equal(t, c.output.code, c.input.rr.Code)
			assert.JSONEq(t, c.output.body, c.input.rr.Body.String())
			assert.Equal(t, c.output.headers, c.input.rr.Header())
			// -> expectations
			mpMock.AssertExpectations(t)
			hdMock.AssertExpectations(t)
		})
	}
}
package handlers

import (
	"api/internal/task"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LNMMusic/optional"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tests
func TestHandlerTask_Get(t *testing.T) {
	type input struct {setW func(w *httptest.ResponseRecorder); setR func(r *http.Request)}
	type output struct {status int; body string}
	type testCase struct {
		title	   string
		input	   input
		output	   output
		setStorage func(mk *task.StorageMock)
	}

	cases := []testCase{
		// succeed cases
		{
			title: "Get a task",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodGet
					r.URL.Path = "/tasks/{id}"
					// context (to get route params from path with chi)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add("id", "1")
					*r = *r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				},
			},
			output: output{
				status: http.StatusOK,
				body: `{
					"message": "succeed to get task",
					"data": {
						"id": "1",
						"title": "title",
						"description": "description",
						"completed": false
					}
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.
					On("Get", mock.Anything).
					Return(&task.Task{
						ID: optional.Some("1"),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(false),
					}, nil)
			},
		},

		// failed cases
		{
			title: "Failed to get a task: not found",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodGet
					r.URL.Path = "/tasks/1"
					// context (to get route params from path with chi)
					chiCtx := chi.NewRouteContext()
					chiCtx.URLParams.Add("id", "1")
					*r = *r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))
				},
			},
			output: output{
				status: http.StatusNotFound,
				body: `{
					"data": null,
					"message": "failed to get task: not found"
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.
					On("Get", mock.Anything).
					Return(&task.Task{}, task.ErrStorageNotFound)
			},
		},
		{
			title: "Failed to get a task: internal error",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodGet
					r.URL.Path = "/tasks/1"
					// context (to get route params from path with chi)
					ctx := chi.NewRouteContext()
					ctx.URLParams.Add("id", "1")
					*r = *r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))
				},
			},
			output: output{
				status: http.StatusInternalServerError,
				body: `{
					"data": null,
					"message": "internal error"
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.
					On("Get", mock.Anything).
					Return(&task.Task{}, task.ErrStorageInternal)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			st := task.NewStorageMock()
			c.setStorage(st)

			cl := NewTaskController(st)
			hd := cl.Get()

			// act
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
			c.input.setR(r)
			hd(w, r)

			// assert
			assert.Equal(t, c.output.status, w.Code)
			assert.JSONEq(t, c.output.body, w.Body.String())
			st.AssertExpectations(t)
		})
	}
}

func TestHandlerTask_Create(t *testing.T) {
	type input struct {setW func(w *httptest.ResponseRecorder); setR func(r *http.Request)}
	type output struct {status int; body string}
	type testCase struct {
		title	   string
		input	   input
		output	   output
		setStorage func(mk *task.StorageMock)
	}

	cases := []testCase{
		// succeed cases
		{
			title: "Create a task",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodPost
					r.URL.Path = "/tasks"
					body := strings.NewReader(`{
						"title": "title",
						"description": "description",
						"completed": false
					}`)
					r.Body = io.NopCloser(body)
				},
			},
			output: output{
				status: http.StatusCreated,
				body: `{
					"message": "succeed to create task",
					"data": {
						"id": "1",
						"title": "title",
						"description": "description",
						"completed": false
					}
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.SetTask = func(t *task.Task) {
					t.ID = optional.Some("1")
				}
				mk.
					On("Save", &task.Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(false),
					}).
					Return(nil)
			},
		},

		// failed cases
		{
			title: "Failed to create a task: decoder",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodPost
					r.URL.Path = "/tasks"
					body := strings.NewReader(`{wrong decoder}`)
					r.Body = io.NopCloser(body)
				},
			},
			output: output{
				status: http.StatusBadRequest,
				body: `{
					"data": null,
					"message": "failed to create task: invalid request"
				}`,
			},
			setStorage: func(mk *task.StorageMock) {},
		},
		{
			title: "Failed to create a task: validator",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodPost
					r.URL.Path = "/tasks"
					body := strings.NewReader(`{
						"title": null,
						"description": null,
						"completed": null
					}`)
					r.Body = io.NopCloser(body)
				},
			},
			output: output{
				status: http.StatusUnprocessableEntity,
				body: `{
					"data": null,
					"message": "failed to create task: invalid task"
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.
					On("Save", &task.Task{
						ID: optional.None[string](),
						Title: optional.None[string](),
						Description: optional.None[string](),
						Completed: optional.None[bool](),
					}).
					Return(task.ErrStorageInvalid)
			},
		},
		{
			title: "Failed to create a task: internal error",
			input: input{
				setW: func(w *httptest.ResponseRecorder) {},
				setR: func(r *http.Request) {
					// base
					r.Method = http.MethodPost
					r.URL.Path = "/tasks"
					body := strings.NewReader(`{
						"title": "title",
						"description": "description",
						"completed": false
					}`)
					r.Body = io.NopCloser(body)
				},
			},
			output: output{
				status: http.StatusInternalServerError,
				body: `{
					"data": null,
					"message": "internal error"
				}`,
			},
			setStorage: func(mk *task.StorageMock) {
				mk.
					On("Save", &task.Task{
						ID: optional.None[string](),
						Title: optional.Some("title"),
						Description: optional.Some("description"),
						Completed: optional.Some(false),
					}).
					Return(task.ErrStorageInternal)
			},
		},
	}

	// run tests
	for _, c := range cases {
		t.Run(c.title, func(t *testing.T) {
			// arrange
			st := task.NewStorageMock()
			c.setStorage(st)

			cl := NewTaskController(st)
			hd := cl.Create()

			// act
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/tasks", nil)
			c.input.setR(r)
			hd(w, r)

			// assert
			assert.Equal(t, c.output.status, w.Code)
			assert.JSONEq(t, c.output.body, w.Body.String())
			st.AssertExpectations(t)
		})
	}
}
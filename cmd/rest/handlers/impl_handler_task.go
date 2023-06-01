package handlers

import (
	"api/cmd/rest/middlewares/logger"
	"api/cmd/rest/response"
	"api/internal/task"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LNMMusic/optional"

	"github.com/go-chi/chi"
)

func NewTaskController(storage task.Storage) *Task {
	return &Task{storage: storage}
}

// Task is an implementation of the task controller.
type Task struct {
	// storage
	storage task.Storage
}

func (t *Task) Get() http.HandlerFunc {
	type resp struct {
		ID			optional.Option[string]	`json:"id"`
		Title		optional.Option[string]	`json:"title"`
		Description	optional.Option[string]	`json:"description"`
		Completed	optional.Option[bool]	`json:"completed"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// param id
		id := chi.URLParam(r, "id")

		// process
		ts, err := t.storage.Get(id)
		if err != nil {
			switch {
				case errors.Is(err, task.ErrStorageNotFound):
					response.Err(w, http.StatusNotFound, "failed to get task: not found")
				default:
					response.Err(w, http.StatusInternalServerError, "internal error")
			}
			logger.Errors(r, err)

			return
		}

		// response
		response.Ok(w, http.StatusOK, "succeed to get task", resp{
			ID: 		 ts.ID,
			Title: 		 ts.Title,
			Description: ts.Description,
			Completed: 	 ts.Completed,
		})
	}
}

func (t *Task) Create() http.HandlerFunc {
	type request struct {
		Title 		optional.Option[string] `json:"title"`
		Description optional.Option[string] `json:"description"`
		Completed 	optional.Option[bool]	`json:"completed"`
	}

	type resp struct {
		ID			optional.Option[string]	`json:"id"`
		Title		optional.Option[string]	`json:"title"`
		Description	optional.Option[string]	`json:"description"`
		Completed	optional.Option[bool]	`json:"completed"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var req request
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			response.Err(w, http.StatusBadRequest, "failed to create task: invalid request")
			logger.Errors(r, err)
			return
		}

		// process
		ts := &task.Task{
			ID: 		 optional.None[string](),
			Title: 		 req.Title,
			Description: req.Description,
			Completed: 	 req.Completed,
		}
		err = t.storage.Save(ts)
		if err != nil {
			switch {
				case errors.Is(err, task.ErrStorageInvalid):
					response.Err(w, http.StatusUnprocessableEntity, "failed to create task: invalid task")
				default:
					response.Err(w, http.StatusInternalServerError, "internal error")
			}
			logger.Errors(r, err)

			return
		}

		// response
		response.Ok(w, http.StatusCreated, "succeed to create task", resp{
			ID: 		 ts.ID,
			Title: 		 ts.Title,
			Description: ts.Description,
			Completed: 	 ts.Completed,
		})
	}
}
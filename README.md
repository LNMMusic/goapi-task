# Go API

This project is a simple API built in Go programming language that allows users to retrieve and save tasks. It follows an interface-oriented architecture and utilizes the `chi` package.

## Main Interfaces

The main interfaces of the task are defined in the `task` package. These interfaces include:

### Task Struct

The `Task` struct represents a task and contains the following optional fields:

```go
type Task struct {
	ID          optional.Option[string]
	Title       optional.Option[string]
	Description optional.Option[string]
	Completed   optional.Option[bool]
}
```

### Storage Interface

The `Storage` interface defines the basic methods for task storage:

```go
type Storage interface {
	Get(id string) (ts *Task, err error)
	Save(task *Task) (err error)
}
```

The following error variables are also defined:

```go
var (
	ErrStorageInternal = errors.New("storage internal error")
	ErrStorageNotFound = errors.New("storage task not found")
	ErrStorageInvalid  = errors.New("storage invalid task")
)
```

### Validator Interface

The `Validator` interface defines the basic methods for task validation:

```go
type Validator interface {
	Validate(task *Task) (err error)
}
```

The following error variables are also defined:

```go
var (
	ErrValidatorInternal      = errors.New("validator internal error")
	ErrValidatorFieldRequired = errors.New("validator field required")
	ErrValidatorFieldEmpty    = errors.New("validator field empty")
	ErrValidatorFieldQuality  = errors.New("validator field quality")
)
```

## Implementations

The project provides the following implementations for storage and validator:

### Storage Implementations

1. **Local Storage**: Implements the `Storage` interface and utilizes a model that allows null values. It uses the `optional` package to handle null values. The local storage implementation is defined in the `task` package.

```go
package task

type LocalStorage struct {
	db []*Task
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		db: make([]*Task, 0),
	}
}

func (s *LocalStorage) Get(id string) (*Task, error) {
	// Implementation for retrieving a task by ID from local storage
}

func (s *LocalStorage) Save(task *Task) error {
	// Implementation for saving a task to local storage
}
```

2. **MySQL Storage**: Implements the `Storage` interface using MySQL database. It utilizes a DTO (Data Transfer Object) with SQL null or value types.

```go
package task

import (
	"database/sql"
	"fmt"

	"github.com/LNMMusic/optional"
	"github.com/google/uuid"
)

// constructor
func NewStorageMySQL(db *sql.DB, vl Validator) *StorageMySQL {
	return &StorageMySQL{db: db, vl: vl}
}

// StorageMySQL is an implementation with MySQL of the Storage interface.
const (
	QueryGetTask = `SELECT id, title, description, completed FROM tasks WHERE id = ?`
	QuerySaveTask = `INSERT INTO tasks (id, title, description, completed) VALUES (?, ?, ?, ?)`
)

// TaskMySQL is the MySQL representation of a task. (internal Data Transfer Object)
type TaskMySQL struct {
	ID          sql.NullString
	Title       sql.NullString
	Description sql.NullString
	Completed   sql.NullBool
}

type StorageMySQL struct {
	// db is the database connection.
	db *sql.DB
	// vl is the task validator.
	vl Validator
}

// Get returns the task with the given id.
func (s *StorageMySQL) Get(id string) (ts *Task, err error) {
	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(QueryGetTask)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "prepare")
		return
	}
	defer stmt.Close()

	// execute statement
	var taskMySQL TaskMySQL
	err = stmt.QueryRow(id).Scan(&taskMySQL.ID, &taskMySQL.Title, &taskMySQL.Description, &taskMySQL.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("%w: %s", ErrStorageNotFound, "query row")
			return
		}
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "query row")
		return
	}

	// serialize
	ts = &Task{}
	if taskMySQL.ID.Valid {
		ts.ID = optional.Some(taskMySQL.ID.String)
	}
	if taskMySQL.Title.Valid {
		ts.Title = optional.Some(taskMySQL.Title.String)
	}
	if taskMySQL.Description.Valid {
		ts.Description = optional.Some(taskMySQL.Description.String)
	}
	if taskMySQL.Completed.Valid {
		ts.Completed = optional.Some(taskMySQL.Completed.Bool)
	}

	return
}
```

In this case, the implementation of mysql needs to use some specific types such as `sql.NullString` to allow compatibility with null values from go to mysql. For this reason, during the query of `SELECT` a `serialization` is apply from the struct `TaskMySQL` (using sql types) to `Task` (using optional type).

In the MySQL storage implementation, the `StorageMySQL` struct implements the `Storage` interface. It has a database connection (`db`) and a task validator (`vl`) as its fields. The `Get` method retrieves a task by its ID from the MySQL database, and the `Save` method saves a task to the MySQL storage.

Make sure to update the MySQL connection details and the query strings according to your specific MySQL setup.



### Validator Implementation

The project provides a validator implementation called **Local Validator**, which implements the `Validator` interface. This implementation performs local validation of tasks.

```go
package task

type LocalValidator struct {
}

func NewValidatorLocal() *LocalValidator {
	return &LocalValidator{}
}

func (v *LocalValidator) Validate(task *Task) error {
	// Implementation for validating a task locally
}
```

The `LocalValidator` struct implements the `Validator` interface and has a `Validate` method to validate tasks.

## Initialization

The application is initialized in the `application` package. It sets up the necessary dependencies and registers the routes using the `chi` router.

```go
package application

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ...

type App struct {
	config *Config
	router chi.Router
}

func (a *App) Dependencies() error {
	// Initialize dependencies (based on config)
	storage := task.NewLocalStorage()
	validator := task.NewLocalValidator()
	// ...

	// Register routes


	// ...

	return nil
}

func (a *App) Run() error {
	// Start the application
	err := http.ListenAndServe(":8080", a.router)
	if err != nil {
		return errors.New("failed to start the server: " + err.Error())
	}

	return nil
}
```

## Example Usage

```go
package main

import (
	"api/cmd/rest/application"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// env
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	
	// app
	config := application.NewConfigDefault()
	router := chi.NewRouter()

	app := application.NewApp(config, router)
	if err := app.Dependencies(); err != nil {
		panic(err)
	}

	// run
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

The API exposes the storage service on a Chi server, and the following routes are registered:

- `GET /ping`: Health check endpoint.
- `GET /tasks/{id}`: Retrieves a task by its ID.
- `POST /tasks`: Creates a new task.

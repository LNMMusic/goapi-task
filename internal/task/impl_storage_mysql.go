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
	ID 			sql.NullString
	Title 		sql.NullString
	Description sql.NullString
	Completed 	sql.NullBool
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

// Save saves the given task.
func (s *StorageMySQL) Save(task *Task) (err error) {
	// validate
	err = s.vl.Validate(task)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInvalid, "validate")
		return
	}
	
	// deserialize
	var taskMySQL TaskMySQL
	if task.Title.IsSome() {
		taskMySQL.Title.String, _ = task.Title.Unwrap()
		taskMySQL.Title.Valid = true
	}
	if task.Description.IsSome() {
		taskMySQL.Description.String, _ = task.Description.Unwrap()
		taskMySQL.Description.Valid = true
	}
	if task.Completed.IsSome() {
		taskMySQL.Completed.Bool, _ = task.Completed.Unwrap()
		taskMySQL.Completed.Valid = true
	}
		
	// default values
	taskMySQL.ID.String = uuid.New().String()
	taskMySQL.ID.Valid = true
	
	// prepare transaction
	var tx *sql.Tx
	tx, err = s.db.Begin()
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "begin")
		return
	}
	defer func () {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(QuerySaveTask)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "prepare")
		return
	}
	defer stmt.Close()

	// execute statement
	var result sql.Result
	result, err = stmt.Exec(taskMySQL.ID, taskMySQL.Title, taskMySQL.Description, taskMySQL.Completed)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "exec")
		return
	}

	// check result
	var rowsAffected int64
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "result rows affected")
		return
	}

	// check rows affected
	if rowsAffected != 1 {
		err = fmt.Errorf("%w: %s", ErrStorageInternal, "rows affected")
		return
	}

	return
}

package mapper

import (
	"database/sql"
	"errors"
	"fmt"
)

// NewProfileMapperMySQL returns a new instance of the MySQL mapper
func NewProfileMapperMySQL(db *sql.DB) *ProfileMapperMySQL {
	return &ProfileMapperMySQL{db: db}
}

// MapperMySQL is the MySQL implementation of the mapper interface
type ProfileMapperMySQL struct {
	// db is the database connection
	db *sql.DB
}

func (impl *ProfileMapperMySQL) MapProfile(userId string) (profileId string, err error) {
	// query
	query := "SELECT profile_id FROM profiles WHERE user_id = ?"

	// prepare
	var stmt *sql.Stmt
	stmt, err = impl.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrProfileMapperInternal, err.Error())
		return
	}

	// execute
	row := stmt.QueryRow(userId)
	if row.Err() != nil {
		switch {
		case errors.Is(row.Err(), sql.ErrNoRows):
			err = fmt.Errorf("%w. %s", ErrProfileMapperNotFound, row.Err().Error())
		default:
			err = fmt.Errorf("%w. %s", ErrProfileMapperInternal, row.Err().Error())
		}
		return
	}

	// scan
	err = row.Scan(&profileId)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrProfileMapperInternal, err.Error())
		return
	}

	return
}
package profiles

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/LNMMusic/optional"
	"github.com/go-sql-driver/mysql"
)

// ProfileMySQL is a dto for profiles in MySQL
type ProfileMySQL struct {
	ID 	    sql.NullString
	UserID  sql.NullString
	Name    sql.NullString
	Email   sql.NullString
	Phone   sql.NullString
	Address sql.NullString
}

func NewImplStorageMySQL(db *sql.DB) (s *ImplStorageMySQL) {
	s = &ImplStorageMySQL{
		db: db,
	}
	return
}

// ImplStorageMySQL is the implementation of the Storage interface for MySQL
type ImplStorageMySQL struct {
	// db is the database connection
	db *sql.DB
}

// GetProfileByUserId returns a profile by its userId
func (s *ImplStorageMySQL) GetProfileByUserId(userId string) (pf *Profile, err error) {
	// query
	query := "SELECT id, user_id, name, email, phone, address FROM profiles WHERE user_id = ?"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	// execute query
	row := stmt.QueryRow(userId)
	if row.Err() != nil {
		switch {
		case errors.Is(row.Err(), sql.ErrNoRows):
			err = fmt.Errorf("%w. %s", ErrStorageNotFound, row.Err().Error())
		default:
			err = fmt.Errorf("%w. %s", ErrStorageInternal, row.Err().Error())
		}
		return
	}

	// scan row
	var pfMySQL ProfileMySQL
	err = row.Scan(&pfMySQL.ID, &pfMySQL.UserID, &pfMySQL.Name, &pfMySQL.Email, &pfMySQL.Phone, &pfMySQL.Address)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	// serialize ProfileMySQL to Profile
	pf = new(Profile)
	if pfMySQL.ID.Valid {
		pf.ID = optional.Some(pfMySQL.ID.String)
	}
	if pfMySQL.UserID.Valid {
		pf.UserID = optional.Some(pfMySQL.UserID.String)
	}
	if pfMySQL.Name.Valid {
		pf.Name = optional.Some(pfMySQL.Name.String)
	}
	if pfMySQL.Email.Valid {
		pf.Email = optional.Some(pfMySQL.Email.String)
	}
	if pfMySQL.Phone.Valid {
		pf.Phone = optional.Some(pfMySQL.Phone.String)
	}
	if pfMySQL.Address.Valid {
		pf.Address = optional.Some(pfMySQL.Address.String)
	}

	return
}

// ActivateProfile
func (s *ImplStorageMySQL) ActivateProfile(pf *Profile) (err error) {
	// deserialize Profile to ProfileMySQL
	var pfMySQL ProfileMySQL
	if pf.ID.IsSome() {
		pfMySQL.ID.String, _ = pf.ID.Unwrap()
		pfMySQL.ID.Valid = true
	}
	if pf.UserID.IsSome() {
		pfMySQL.UserID.String, _ = pf.UserID.Unwrap()
		pfMySQL.UserID.Valid = true
	}
	if pf.Name.IsSome() {
		pfMySQL.Name.String, _ = pf.Name.Unwrap()
		pfMySQL.Name.Valid = true
	}
	if pf.Email.IsSome() {
		pfMySQL.Email.String, _ = pf.Email.Unwrap()
		pfMySQL.Email.Valid = true
	}
	if pf.Phone.IsSome() {
		pfMySQL.Phone.String, _ = pf.Phone.Unwrap()
		pfMySQL.Phone.Valid = true
	}
	if pf.Address.IsSome() {
		pfMySQL.Address.String, _ = pf.Address.Unwrap()
		pfMySQL.Address.Valid = true
	}

	// query
	query := "INSERT INTO profiles (id, user_id, name, email, phone, address) VALUES (?, ?, ?, ?, ?, ?)"

	// prepare statement
	var stmt *sql.Stmt
	stmt, err = s.db.Prepare(query)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	// execute query
	var result sql.Result
	result, err = stmt.Exec(pfMySQL.ID, pfMySQL.UserID, pfMySQL.Name, pfMySQL.Email, pfMySQL.Phone, pfMySQL.Address)
	if err != nil {
		errMySQL, ok := err.(*mysql.MySQLError)
		if ok {
			switch errMySQL.Number {
			case 1062:
				err = fmt.Errorf("%w. %s", ErrStorageNotUnique, err.Error())
			default:
				err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
			}
			return
		}

		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	// check affected rows
	var affectedRows int64
	affectedRows, err = result.RowsAffected()
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	if affectedRows != 1 {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, "rows affected != 1")
		return
	}

	return
}
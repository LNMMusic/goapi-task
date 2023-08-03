package transactioner

import (
	"database/sql"
	"fmt"
)

// operation is a function that can be run in a transaction
type operation func() (err error)

// NewImplTransactionerDefault creates a new default implementation of Transactioner
func NewImplTransactionerDefault(db *sql.DB) (impl *ImplTransactionerDefault) {
	impl = &ImplTransactionerDefault{
		db: db,
	}
	return
}

// ImplTransactionerDefault is the default implementation of Transactioner
type ImplTransactionerDefault struct {
	db *sql.DB
}

func (impl *ImplTransactionerDefault) Do(op operation) (err error) {
	// begin transaction
	var tx *sql.Tx
	tx, err = impl.db.Begin()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrTransactionBegin, err)
		return
	}
	// defer rollback/commit
	defer func() {
		// rollback transaction
		if err != nil {
			e := tx.Rollback()
			if e != nil {
				err = fmt.Errorf("%w. %v", ErrTransactionRollback, e)
				return
			}

			return
		}

		// commit transaction
		err = tx.Commit()
		if err != nil {
			err = fmt.Errorf("%w. %v", ErrTransactionCommit, err)
			return
		}
	}()

	// run operation
	err = op()
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrTransactionOperation, err)
		return
	}
	
	return
}
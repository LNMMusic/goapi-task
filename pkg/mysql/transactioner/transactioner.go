package transactioner

import "errors"

// Transactioner is an interface for a mysql transaction
type Transactioner interface {
	// Do runs the operation in a transaction
	// - Success: transaction is committed
	// - Failure: transaction is rolled back
	Do(op operation) (err error)
}
var (
	// ErrTransactionBegin is returned when a transaction cannot be started
	ErrTransactionBegin = errors.New("transactioner: cannot begin transaction")
	// ErrTransactionOp is returned when an operation fails
	ErrTransactionOperation = errors.New("transactioner: operation failed")
	// ErrTransactionCommit is returned when a transaction cannot be committed
	ErrTransactionCommit = errors.New("transactioner: cannot commit transaction")
	// ErrTransactionRollback is returned when a transaction cannot be rolled back
	ErrTransactionRollback = errors.New("transactioner: cannot rollback transaction")
)
package profiles

import (
	"api/pkg/mysql/transactioner"
	"errors"
	"fmt"
)

func NewImplStorageMySQLTx(st Storage, tr transactioner.Transactioner) *ImplStorageMySQLTx {
	return &ImplStorageMySQLTx{
		st: st,
		tr: tr,
	}
}

// ImplStorageMySQLTx is the implementation of the Storage interface for MySQL
// - Wraps the ImplStorageMySQL with a transaction
type ImplStorageMySQLTx struct {
	// st is the storage implementation (to be wrapped)
	st Storage
	// transactioner is the transactioner implementation
	tr transactioner.Transactioner
}

// GetProfileById returns a profile by its userId
func (s *ImplStorageMySQLTx) GetProfileById(id string) (pf *Profile, err error) {
	// run operation
	e := s.tr.Do(func() (e error) {
		// get base values from storage (wrapping process)
		pf, err = s.st.GetProfileById(id)
		if err != nil {
			e = err
		}
		return
	})
	if e != nil {
		switch {
		case errors.Is(e, transactioner.ErrTransactionOperation):
			return
		default:
			err = fmt.Errorf("%w. %s", ErrStorageInternal, e.Error())
		}
		return
	}

	return
}

// ActivateProfile
func (s *ImplStorageMySQLTx) ActivateProfile(pf *Profile) (err error) {
	// run operation
	e := s.tr.Do(func() (e error) {
		// get base values from storage (wrapping process)
		err = s.st.ActivateProfile(pf)
		if err != nil {
			e = err
		}
		return
	})
	if e != nil {
		switch {
		case errors.Is(e, transactioner.ErrTransactionOperation):
			return
		default:
			err = fmt.Errorf("%w. %s", ErrStorageInternal, e.Error())
		}
		return
	}

	return
}
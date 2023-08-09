package storage

import (
	"api/internal/profiles"
	"api/pkg/mysql/transactioner"
	"errors"
	"fmt"
)

func NewImplProfilesStorageMySQLTx(st ProfilesStorage, tr transactioner.Transactioner) *ImplProfilesStorageMySQLTx {
	return &ImplProfilesStorageMySQLTx{
		st: st,
		tr: tr,
	}
}

// ImplProfilesStorageMySQLTx is the implementation of the Storage interface for MySQL
// - Wraps the ImplProfilesStorageMySQL with a transaction
type ImplProfilesStorageMySQLTx struct {
	// st is the storage implementation (to be wrapped)
	st ProfilesStorage
	// transactioner is the transactioner implementation
	tr transactioner.Transactioner
}

// GetProfileById returns a profile by its userId
func (s *ImplProfilesStorageMySQLTx) GetProfileById(id string) (pf *profiles.Profile, err error) {
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
func (s *ImplProfilesStorageMySQLTx) ActivateProfile(pf *profiles.Profile) (err error) {
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
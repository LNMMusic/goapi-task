package profiles

import (
	"api/pkg/mysql/transactioner"
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

// GetProfileByID returns a profile by its ID
func (s *ImplStorageMySQLTx) GetProfileByID(id string) (pf *Profile, err error) {
	// run operation
	err = s.tr.Do(func() (e error) {
		// get base values from storage (wrapping process)
		pf, err = s.st.GetProfileByID(id)
		if err != nil {
			e = err
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	return
}

// ActivateProfile
func (s *ImplStorageMySQLTx) ActivateProfile(pf *Profile) (err error) {
	// run operation
	err = s.tr.Do(func() (e error) {
		err = s.st.ActivateProfile(pf)
		if err != nil {
			e = err
		}
		return
	})
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInternal, err.Error())
		return
	}

	return
}
package profiles

import "fmt"

func NewImplStorageValidator(st Storage, vl Validator) *ImplStorageValidator {
	return &ImplStorageValidator{
		st: st,
		vl: vl,
	}
}

// ImplStorageValidator is the implementation of the Storage interface using Validator interface
type ImplStorageValidator struct {
	// st is the storage implementation (to be wrapped)
	st Storage

	// vl is the validator implementation
	vl Validator
}

// GetProfileByUserId returns a profile by its userId
func (impl *ImplStorageValidator) GetProfileByUserId(id string) (pf *Profile, err error) {
	pf, err = impl.st.GetProfileByUserId(id)
	return
}

// ActivateProfile
func (impl *ImplStorageValidator) ActivateProfile(pf *Profile) (err error) {
	// validate profile
	err = impl.vl.Validate(pf)
	if err != nil {
		err = fmt.Errorf("%w. %s", ErrStorageInvalidProfile, err.Error())
		return
	}

	// activate profile
	err = impl.st.ActivateProfile(pf)
	return 
}

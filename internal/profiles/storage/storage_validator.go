package storage

import (
	"api/internal/profiles"
	"api/internal/profiles/validator"
	"fmt"
)

// NewImplProfilesStorageValidator returns a new instance of ImplProfilesStorageValidator
func NewImplProfilesStorageValidator(st ProfilesStorage, vl validator.ProfilesValidator) *ImplProfilesStorageValidator {
	return &ImplProfilesStorageValidator{
		st: st,
		vl: vl,
	}
}

// ImplProfilesStorageValidator is the implementation of the ProfilesStorage interface using validator.ProfilesValidator interface
type ImplProfilesStorageValidator struct {
	// st is the storage implementation (to be wrapped)
	st ProfilesStorage

	// vl is the validator implementation
	vl validator.ProfilesValidator
}

// GetProfileById returns a profile by its userId
func (impl *ImplProfilesStorageValidator) GetProfileById(id string) (pf *profiles.Profile, err error) {
	pf, err = impl.st.GetProfileById(id)
	return
}

// ActivateProfile
func (impl *ImplProfilesStorageValidator) ActivateProfile(pf *profiles.Profile) (err error) {
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

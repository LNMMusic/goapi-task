package profiles

import (
	"errors"

	"github.com/LNMMusic/optional"
)

// Interfaces for profiles
type Profile struct {
	// ProfileID is the unique identifier for the profile
	ID 	    optional.Option[string]
	// UserID is the unique identifier for the user from a Central Authentication Service
	UserID  optional.Option[string]
	// Name is the name of the user
	Name    optional.Option[string]
	// Email is the email of the user
	Email   optional.Option[string]
	// Phone is the phone number of the user
	Phone   optional.Option[string]
	// Address is the address of the user
	Address optional.Option[string]
}

// Storage interface for profiles
type Storage interface {
	// GetProfileByID returns a profile by its ID
	GetProfileByUserId(userId string) (pf *Profile, err error)

	// ActivateProfile
	ActivateProfile(pf *Profile) (err error)
}
var (
	ErrStorageInternal		 = errors.New("storage: internal storage error")
	ErrStorageInvalidProfile = errors.New("storage: invalid profile")
	ErrStorageNotFound		 = errors.New("storage: profile not found")
	ErrStorageNotUnique	     = errors.New("storage: profile not unique")
)

// Validator interface for profiles
type Validator interface {
	// DefaultProfile set default values for a profile
	Default(pf *Profile) (err error)

	// ValidateProfile validates a profile
	Validate(pf *Profile) (err error)
}
var (
	ErrValidatorInternal	   = errors.New("validator: internal validator error")
	ErrValidatorInvalidProfile = errors.New("validator: invalid profile")
)
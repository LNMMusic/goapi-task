package validator

import (
	"api/internal/profiles"
	"errors"
)

// ProfilesValidator interface for profiles
type ProfilesValidator interface {
	// DefaultProfile set default values for a profile
	Default(pf *profiles.Profile) (err error)

	// ValidateProfile validates a profile
	Validate(pf *profiles.Profile) (err error)
}

var (
	ErrValidatorInternal	   = errors.New("validator: internal validator error")
	ErrValidatorInvalidProfile = errors.New("validator: invalid profile")
)
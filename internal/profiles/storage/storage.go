package storage

import (
	"api/internal/profiles"
	"errors"
)

// ProfilesStorage interface for profiles
type ProfilesStorage interface {
	// GetProfileByID returns a profile by its ID
	GetProfileById(id string) (pf *profiles.Profile, err error)

	// ActivateProfile
	ActivateProfile(pf *profiles.Profile) (err error)
}

var (
	ErrStorageInternal		 = errors.New("storage: internal storage error")
	ErrStorageInvalidProfile = errors.New("storage: invalid profile")
	ErrStorageNotFound		 = errors.New("storage: profile not found")
	ErrStorageNotUnique	     = errors.New("storage: profile not unique")
)
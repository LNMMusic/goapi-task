package profiles

import (
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
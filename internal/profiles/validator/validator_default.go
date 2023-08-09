package validator

import (
	"api/internal/profiles"
	"fmt"
	"regexp"
)

type Config struct {
	// regex patterns
	RegexEmail string
	RegexPhone string
}
func NewImplProfilesValidatorDefault(cfg *Config) (impl *ImplProfilesValidatorDefault) {
	// default config
	defaultCfg := &Config{
		RegexEmail: `^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`,
		RegexPhone: `^[0-9]{10}$`,
	}
	if cfg != nil {
		if cfg.RegexEmail != "" {
			defaultCfg.RegexEmail = cfg.RegexEmail
		}
		if cfg.RegexPhone != "" {
			defaultCfg.RegexPhone = cfg.RegexPhone
		}
	}

	// compile regex patterns
	regexEmail, _ := regexp.Compile(defaultCfg.RegexEmail)
	regexPhone, _ := regexp.Compile(defaultCfg.RegexPhone)

	// create implementation
	impl = &ImplProfilesValidatorDefault{
		regexEmail: regexEmail,
		regexPhone: regexPhone,
	}
	return
}

// ImplProfilesValidatorDefault is the default implementation of the Validator interface
type ImplProfilesValidatorDefault struct {
	// regex patterns
	regexEmail *regexp.Regexp
	regexPhone *regexp.Regexp
}

func (impl *ImplProfilesValidatorDefault) Default(pf *profiles.Profile) (err error) {
	// set default values for profile
	return
}

func (impl *ImplProfilesValidatorDefault) Validate(pf *profiles.Profile) (err error) {
	// required fields (not null)
	if !pf.ID.IsSome() {
		err = fmt.Errorf("%w - id field is required", ErrValidatorInvalidProfile)
		return
	}
	if !pf.UserID.IsSome() {
		err = fmt.Errorf("%w - user_id field is required", ErrValidatorInvalidProfile)
		return
	}

	// quality validation
	userID, _ := pf.UserID.Unwrap()
	if userID == "" {
		err = fmt.Errorf("%w - user_id field can not be empty", ErrValidatorInvalidProfile)
		return
	}
	if pf.Name.IsSome() {
		name, _ := pf.Name.Unwrap()
		if len(name) < 3 || len(name) > 50 {
			err = fmt.Errorf("%w - name field must be between 3 and 50 characters", ErrValidatorInvalidProfile)
			return
		}
	}
	if pf.Email.IsSome() {
		email, _ := pf.Email.Unwrap()
		if !impl.regexEmail.MatchString(email) {
			err = fmt.Errorf("%w - email field is invalid", ErrValidatorInvalidProfile)
			return
		}
	}
	if pf.Phone.IsSome() {
		phone, _ := pf.Phone.Unwrap()
		if !impl.regexPhone.MatchString(phone) {
			err = fmt.Errorf("%w - phone field is invalid", ErrValidatorInvalidProfile)
			return
		}
	}
	if pf.Address.IsSome() {
		address, _ := pf.Address.Unwrap()
		if len(address) < 3 || len(address) > 50 {
			err = fmt.Errorf("%w - address field must be between 3 and 50 characters", ErrValidatorInvalidProfile)
			return
		}
	}

	return
}

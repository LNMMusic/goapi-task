package handlers

import (
	"api/internal/profiles"
	"api/internal/profiles/contexter"
	"api/pkg/uuidgenerator"
	"api/pkg/web"
	"errors"
	"net/http"

	"github.com/LNMMusic/optional"
)

func NewProfileController(st profiles.Storage, uuid uuidgenerator.UUIDGenerator) *ProfileController {
	return &ProfileController{st: st, uuid: uuid}
}

type ProfileController struct {
	// storage is the storage interface for profiles
	st profiles.Storage
	// uuid is the uuid generator interface
	uuid uuidgenerator.UUIDGenerator
}

// GetProfileByID returns a profile by its ID
// type RequestGetProfileByID struct {} // no need for a request struct
type ProfileDTO struct {
	UserID optional.Option[string] `json:"user_id"`
	Name   optional.Option[string] `json:"name"`
	Email  optional.Option[string] `json:"email"`
	Phone  optional.Option[string] `json:"phone"`
	Address optional.Option[string] `json:"address"`
}
type ResponseGetProfileByID struct {
	Message string		`json:"message"`
	Data    *ProfileDTO `json:"data"`
	Error	bool		`json:"error"`
}
func (ct *ProfileController) GetProfileById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := r.Context().Value(contexter.KeyProfileId).(string)

		// process
		pf, err := ct.st.GetProfileById(id)
		if err != nil {
			var code int; var body *ResponseGetProfileByID

			switch {
			case errors.Is(err, profiles.ErrStorageNotFound):
				code = http.StatusNotFound
				body = &ResponseGetProfileByID{
					Message: "Profile not found",
					Data:    nil,
					Error:   true,
				}
			default:
				code = http.StatusInternalServerError
				body = &ResponseGetProfileByID{
					Message: "Internal server error",
					Data:    nil,
					Error:   true,
				}
			}

			web.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseGetProfileByID{
			Message: "Success",
			Data: &ProfileDTO{
				UserID: pf.UserID,
				Name:   pf.Name,
				Email:  pf.Email,
				Phone:  pf.Phone,
				Address: pf.Address,
			},
			Error: false,
		}

		web.JSON(w, code, body)
	}
}

// ActivateProfile activates a profile
// type RequestActivateProfile struct {} // no need for a request struct
type ResponseActivateProfile struct {
	Message string		`json:"message"`
	Data    any 		`json:"data"`
	Error	bool		`json:"error"`
}
func (ct *ProfileController) ActivateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// -> user id
		userId := (*r).Header.Get("User-Id")

		// process
		pf := new(profiles.Profile)
		pf.ID = optional.Some(ct.uuid.UUID())
		pf.UserID = optional.Some(userId)

		err := ct.st.ActivateProfile(pf)
		if err != nil {
			var code int; var body *ResponseActivateProfile

			switch {
			case errors.Is(err, profiles.ErrStorageInvalidProfile):
				code = http.StatusUnprocessableEntity
				body = &ResponseActivateProfile{
					Message: "Invalid profile",
					Data:    nil,
					Error:   true,
				}
			case errors.Is(err, profiles.ErrStorageNotUnique):
				code = http.StatusConflict
				body = &ResponseActivateProfile{
					Message: "Profile not unique",
					Data:    nil,
					Error:   true,
				}
			default:
				code = http.StatusInternalServerError
				body = &ResponseActivateProfile{
					Message: "Internal server error",
					Data:    nil,
					Error:   true,
				}
			}

			web.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := &ResponseActivateProfile{
			Message: "Success",
			Data:    nil,
			Error:   false,
		}

		web.JSON(w, code, body)
	}
}
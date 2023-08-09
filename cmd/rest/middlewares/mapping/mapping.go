package mapping

import (
	"api/internal/profiles/contexter"
	"api/internal/profiles/mapper"
	"api/pkg/web"
	"context"
	"errors"
	"net/http"
)

// NewProfileMapping returns a new ProfileMapping
func NewProfileMapping(pm mapper.ProfileMapper) *ProfileMapping {
	return &ProfileMapping{ProfileMapper: pm}
}

// ProfileMapping is the mapping interface for profiles
type ProfileMapping struct {
	// mapper
	ProfileMapper mapper.ProfileMapper
}

// MapProfile is a middleware that maps a profile
type ResponseMapProfile struct {
	Message string `json:"message"`
	Data	any `json:"data"`
	Error	bool `json:"error"`
}
func (mp *ProfileMapping) MapProfile(hd http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user id from headers
		userId := r.Header.Get("User-Id")

		// map profile
		profileId, err := mp.ProfileMapper.MapProfile(userId)
		if err != nil {
			var code int; var body *ResponseMapProfile
			switch {
			case errors.Is(err, mapper.ErrProfileMapperNotFound):
				code = http.StatusUnauthorized
				body = &ResponseMapProfile{
					Message: "Profile not found",
					Data:    nil,
					Error:   true,
				}
			default:
				code = http.StatusInternalServerError
				body = &ResponseMapProfile{
					Message: "Internal server error",
					Data:    nil,
					Error:   true,
				}
			}

			web.JSON(w, code, body)
			return
		}
		
		// set profile in context
		(*r) = *(*r).WithContext(context.WithValue((*r).Context(), contexter.KeyProfileId, profileId))

		// next
		hd.ServeHTTP(w, r)
	})
}
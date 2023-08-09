package mapper

import "errors"

type ProfileMapper interface {
	// MapProfile maps user id to profile id
	MapProfile(userId string) (profileId string, err error)
}

var (
	ErrProfileMapperInternal = errors.New("mapper: internal mapper error")
	ErrProfileMapperNotFound = errors.New("mapper: mapper not found")
)
package uuidgenerator

import "errors"

type UUIDGenerator interface {
	// GenerateUUID generates a UUID
	UUID() (id string)
}
var (
	ErrUUIDGeneratorInternal = errors.New("uuidgenerator: internal uuidgenerator error")
)
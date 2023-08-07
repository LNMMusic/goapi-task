package uuidgenerator

import "github.com/google/uuid"

// NewUUIDGeneratorGoogle returns a new UUIDGeneratorGoogle
func NewUUIDGeneratorGoogle() (ug *ImplUUIDGeneratorGoogle) {
	ug = &ImplUUIDGeneratorGoogle{}
	return
}

// ImplUUIDGeneratorGoogle is the implementation of UUIDGenerator using Google's UUID generator
type ImplUUIDGeneratorGoogle struct{}

func (ug *ImplUUIDGeneratorGoogle) UUID() (id string) {
	id = uuid.New().String()
	return
}
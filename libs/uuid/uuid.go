package uuid

import (
	uuidx "github.com/google/uuid"
)

// UUID is the alias wrapper google uuid lib.
type UUID uuidx.UUID

// New creates a new random UUID or panics.
func New() UUID {
	return UUID(uuidx.New())
}

// NewString creates a new random UUID.
// Returns it as a string or panics.
func NewString() string {
	return uuidx.NewString()
}

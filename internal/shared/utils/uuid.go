package utils

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4.
func GenerateUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID.
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// MustGenerateUUID generates a new UUID v4 and panics if it fails.
// This should only be used when UUID generation cannot fail
// (which is extremely unlikely with v4 UUIDs).
func MustGenerateUUID() string {
	id := GenerateUUID()
	if id == "" {
		panic("failed to generate UUID")
	}
	return id
}

// ParseUUID parses a UUID string into a UUID object.
// Returns an error if the string is not a valid UUID.
func ParseUUID(u string) (uuid.UUID, error) {
	return uuid.Parse(u)
}

// UUIDToBytes converts a UUID to a byte slice.
// This is useful when working with databases that store UUIDs as bytes.
func UUIDToBytes(u uuid.UUID) []byte {
	return u[:]
}

// BytesToUUID converts a byte slice to a UUID.
// Returns an error if the byte slice is not a valid UUID.
func BytesToUUID(b []byte) (uuid.UUID, error) {
	var u uuid.UUID
	if len(b) != 16 {
		return uuid.Nil, fmt.Errorf("invalid UUID length: %d", len(b))
	}
	copy(u[:], b)
	return u, nil
}

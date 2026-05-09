package tests

import "github.com/google/uuid"

//go:fix inline
func ptr[T any](v T) *T {
	return new(v)
}

func openapiUUID(id uuid.UUID) uuid.UUID { return id }

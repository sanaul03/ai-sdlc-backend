package vehicle

import "errors"

// Sentinel errors for the vehicle domain.
var (
	ErrNotFound     = errors.New("vehicle not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)

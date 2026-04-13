package fleet

import "errors"

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("fleet: record not found")

// ErrConflict is returned when an operation cannot proceed due to a data conflict
// (e.g., deleting a car group that still has active vehicles).
var ErrConflict = errors.New("fleet: operation conflicts with existing data")

// ErrValidation is returned when input data fails validation rules.
var ErrValidation = errors.New("fleet: validation error")

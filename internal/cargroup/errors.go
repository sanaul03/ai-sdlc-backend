package cargroup

import "errors"

// Sentinel errors for the cargroup domain.
var (
	ErrNotFound      = errors.New("car group not found")
	ErrNameTaken     = errors.New("car group name already in use")
	ErrHasVehicles   = errors.New("car group has active vehicles and cannot be deleted")
	ErrInvalidInput  = errors.New("invalid input")
)

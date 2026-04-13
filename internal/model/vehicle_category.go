package model

import (
	"time"
)

// VehicleCategory represents a top-level grouping for vehicle types (e.g. Bus, Truck).
type VehicleCategory struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateVehicleCategoryRequest holds fields required to create a vehicle category.
type CreateVehicleCategoryRequest struct {
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
}

// UpdateVehicleCategoryRequest holds fields allowed to be updated on a vehicle category.
type UpdateVehicleCategoryRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

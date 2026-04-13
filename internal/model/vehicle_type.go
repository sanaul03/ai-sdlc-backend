package model

import (
	"time"
)

// VehicleType represents a specific vehicle specification within a category.
type VehicleType struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Capacity    *int      `json:"capacity,omitempty"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateVehicleTypeRequest holds fields required to create a vehicle type.
type CreateVehicleTypeRequest struct {
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Capacity    *int    `json:"capacity,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateVehicleTypeRequest holds fields allowed to be updated on a vehicle type.
type UpdateVehicleTypeRequest struct {
	Name        *string `json:"name,omitempty"`
	Capacity    *int    `json:"capacity,omitempty"`
	Description *string `json:"description,omitempty"`
}

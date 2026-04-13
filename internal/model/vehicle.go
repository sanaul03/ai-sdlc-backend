package model

import (
	"time"
)

// VehicleStatus enumerates the allowed values for Vehicle.Status.
const (
	VehicleStatusActive        = "active"
	VehicleStatusInactive      = "inactive"
	VehicleStatusMaintenance   = "maintenance"
	VehicleStatusDecommissioned = "decommissioned"
)

// Vehicle represents an individual vehicle in the fleet.
type Vehicle struct {
	ID                 int64      `json:"id"`
	CompanyID          int64      `json:"company_id"`
	DepotID            *int64     `json:"depot_id,omitempty"`
	VehicleTypeID      int64      `json:"vehicle_type_id"`
	RegistrationNumber string     `json:"registration_number"`
	ChassisNumber      *string    `json:"chassis_number,omitempty"`
	EngineNumber       *string    `json:"engine_number,omitempty"`
	ManufactureYear    *int16     `json:"manufacture_year,omitempty"`
	Color              *string    `json:"color,omitempty"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

// CreateVehicleRequest holds fields required to create a vehicle.
type CreateVehicleRequest struct {
	CompanyID          int64   `json:"company_id"`
	DepotID            *int64  `json:"depot_id,omitempty"`
	VehicleTypeID      int64   `json:"vehicle_type_id"`
	RegistrationNumber string  `json:"registration_number"`
	ChassisNumber      *string `json:"chassis_number,omitempty"`
	EngineNumber       *string `json:"engine_number,omitempty"`
	ManufactureYear    *int16  `json:"manufacture_year,omitempty"`
	Color              *string `json:"color,omitempty"`
}

// UpdateVehicleRequest holds fields allowed to be updated on a vehicle.
type UpdateVehicleRequest struct {
	DepotID         *int64  `json:"depot_id,omitempty"`
	VehicleTypeID   *int64  `json:"vehicle_type_id,omitempty"`
	ChassisNumber   *string `json:"chassis_number,omitempty"`
	EngineNumber    *string `json:"engine_number,omitempty"`
	ManufactureYear *int16  `json:"manufacture_year,omitempty"`
	Color           *string `json:"color,omitempty"`
	Status          *string `json:"status,omitempty"`
}

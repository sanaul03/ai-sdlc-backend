package model

import (
	"time"
)

// Depot represents a physical base or garage within a company's fleet.
type Depot struct {
	ID        int64      `json:"id"`
	CompanyID int64      `json:"company_id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Address   *string    `json:"address,omitempty"`
	Latitude  *float64   `json:"latitude,omitempty"`
	Longitude *float64   `json:"longitude,omitempty"`
	Status    int16      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CreateDepotRequest holds fields required to create a depot.
type CreateDepotRequest struct {
	CompanyID int64    `json:"company_id"`
	Name      string   `json:"name"`
	Code      string   `json:"code"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

// UpdateDepotRequest holds fields allowed to be updated on a depot.
type UpdateDepotRequest struct {
	Name      *string  `json:"name,omitempty"`
	Address   *string  `json:"address,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Status    *int16   `json:"status,omitempty"`
}

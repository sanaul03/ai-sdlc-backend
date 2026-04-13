// Package fleet implements the Fleet Structure domain: car groups and vehicles.
package fleet

import (
	"time"

	"github.com/google/uuid"
)

// CarGroup represents a classification hierarchy entry for rental vehicles
// (e.g., Economy Sedan, Luxury SUV).
type CarGroup struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	SizeCategory *string    `json:"size_category,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	CreatedBy    string     `json:"created_by"`
	UpdatedBy    string     `json:"updated_by"`
	Deleted      bool       `json:"deleted"`
}

// CreateCarGroupInput carries the data required to create a new car group.
type CreateCarGroupInput struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
	CreatedBy    string  `json:"-"`
}

// UpdateCarGroupInput carries the fields that may be changed on an existing car group.
type UpdateCarGroupInput struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
	UpdatedBy    string  `json:"-"`
}

// ListCarGroupsFilter defines optional filters for listing car groups.
type ListCarGroupsFilter struct {
	// Q is a case-insensitive substring search on the name field.
	Q *string
	// IncludeDeleted controls whether soft-deleted records are included.
	IncludeDeleted bool
}

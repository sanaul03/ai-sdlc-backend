package cargroup

import (
	"time"

	"github.com/google/uuid"
)

// CarGroup represents the car_groups database row.
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

// CreateRequest holds the fields required to create a new car group.
type CreateRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
}

// UpdateRequest holds the fields that may be updated on a car group.
type UpdateRequest struct {
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
}

// ListFilter carries optional filters for the list endpoint.
type ListFilter struct {
	Q       string
	Deleted bool
}

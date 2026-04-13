package model

import (
	"time"
)

// Company represents a fleet-owning organisation.
type Company struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Address   *string    `json:"address,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Status    int16      `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CreateCompanyRequest holds fields required to create a company.
type CreateCompanyRequest struct {
	Name    string  `json:"name"`
	Code    string  `json:"code"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
}

// UpdateCompanyRequest holds fields allowed to be updated on a company.
type UpdateCompanyRequest struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	Status  *int16  `json:"status,omitempty"`
}

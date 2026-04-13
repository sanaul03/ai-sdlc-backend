package vehicle

import (
	"time"

	"github.com/google/uuid"
)

// Vehicle represents the vehicles database row.
type Vehicle struct {
	ID                     uuid.UUID  `json:"id"`
	CarGroupID             uuid.UUID  `json:"car_group_id"`
	BranchID               uuid.UUID  `json:"branch_id"`
	VIN                    string     `json:"vin"`
	LicencePlate           string     `json:"licence_plate"`
	Brand                  string     `json:"brand"`
	Model                  string     `json:"model"`
	Year                   int        `json:"year"`
	Colour                 *string    `json:"colour,omitempty"`
	FuelType               string     `json:"fuel_type"`
	TransmissionType       string     `json:"transmission_type"`
	CurrentMileage         int        `json:"current_mileage"`
	Status                 string     `json:"status"`
	Designation            string     `json:"designation"`
	AcquisitionDate        time.Time  `json:"acquisition_date"`
	OwnershipType          string     `json:"ownership_type"`
	LeaseDetails           *string    `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string    `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time `json:"next_inspection_due_date,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
	CreatedBy              string     `json:"created_by"`
	UpdatedBy              string     `json:"updated_by"`
	Deleted                bool       `json:"deleted"`
	ExpiryWarning          bool       `json:"expiry_warning"`
}

// CreateRequest holds the fields used to create a new vehicle.
type CreateRequest struct {
	CarGroupID             uuid.UUID  `json:"car_group_id"`
	BranchID               uuid.UUID  `json:"branch_id"`
	VIN                    string     `json:"vin"`
	LicencePlate           string     `json:"licence_plate"`
	Brand                  string     `json:"brand"`
	Model                  string     `json:"model"`
	Year                   int        `json:"year"`
	Colour                 *string    `json:"colour,omitempty"`
	FuelType               string     `json:"fuel_type"`
	TransmissionType       string     `json:"transmission_type"`
	CurrentMileage         int        `json:"current_mileage"`
	Status                 string     `json:"status"`
	Designation            string     `json:"designation"`
	AcquisitionDate        time.Time  `json:"acquisition_date"`
	OwnershipType          string     `json:"ownership_type"`
	LeaseDetails           *string    `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string    `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time `json:"next_inspection_due_date,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
}

// UpdateRequest holds the updatable fields for a vehicle (status excluded).
type UpdateRequest struct {
	CarGroupID             *uuid.UUID `json:"car_group_id,omitempty"`
	BranchID               *uuid.UUID `json:"branch_id,omitempty"`
	VIN                    *string    `json:"vin,omitempty"`
	LicencePlate           *string    `json:"licence_plate,omitempty"`
	Brand                  *string    `json:"brand,omitempty"`
	Model                  *string    `json:"model,omitempty"`
	Year                   *int       `json:"year,omitempty"`
	Colour                 *string    `json:"colour,omitempty"`
	FuelType               *string    `json:"fuel_type,omitempty"`
	TransmissionType       *string    `json:"transmission_type,omitempty"`
	CurrentMileage         *int       `json:"current_mileage,omitempty"`
	Designation            *string    `json:"designation,omitempty"`
	AcquisitionDate        *time.Time `json:"acquisition_date,omitempty"`
	OwnershipType          *string    `json:"ownership_type,omitempty"`
	LeaseDetails           *string    `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string    `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time `json:"next_inspection_due_date,omitempty"`
	Notes                  *string    `json:"notes,omitempty"`
}

// DesignationUpdateRequest holds the new designation value.
type DesignationUpdateRequest struct {
	Designation string `json:"designation"`
}

// ListFilter carries optional filters for the vehicle list endpoint.
type ListFilter struct {
	CarGroupID       *uuid.UUID
	BranchID         *uuid.UUID
	Status           string
	Designation      string
	FuelType         string
	TransmissionType string
	ExpiryWarning    bool
	Page             int
	PageSize         int
}

// Page holds paginated vehicle results.
type Page struct {
	Items    []Vehicle `json:"items"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}

// Allowed value sets used in validation.
var (
	AllowedFuelTypes         = map[string]struct{}{"petrol": {}, "diesel": {}, "electric": {}, "hybrid": {}}
	AllowedTransmissionTypes = map[string]struct{}{"manual": {}, "automatic": {}}
	AllowedDesignations      = map[string]struct{}{"rental_only": {}, "sales_only": {}, "shared": {}}
	AllowedStatuses          = map[string]struct{}{
		"available":         {},
		"on_rent":           {},
		"needs_cleaning":    {},
		"needs_inspection":  {},
		"under_maintenance": {},
		"unavailable":       {},
		"decommissioned":    {},
	}
)

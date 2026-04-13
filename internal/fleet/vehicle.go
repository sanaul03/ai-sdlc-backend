package fleet

import (
	"time"

	"github.com/google/uuid"
)

// VehicleStatus enumerates the allowed operational statuses for a vehicle.
type VehicleStatus string

const (
	VehicleStatusAvailable        VehicleStatus = "available"
	VehicleStatusOnRent           VehicleStatus = "on_rent"
	VehicleStatusNeedsCleaning    VehicleStatus = "needs_cleaning"
	VehicleStatusNeedsInspection  VehicleStatus = "needs_inspection"
	VehicleStatusUnderMaintenance VehicleStatus = "under_maintenance"
	VehicleStatusUnavailable      VehicleStatus = "unavailable"
	VehicleStatusDecommissioned   VehicleStatus = "decommissioned"
)

// VehicleDesignation enumerates the allowed availability designations.
type VehicleDesignation string

const (
	DesignationRentalOnly VehicleDesignation = "rental_only"
	DesignationSalesOnly  VehicleDesignation = "sales_only"
	DesignationShared     VehicleDesignation = "shared"
)

// FuelType enumerates the allowed fuel types.
type FuelType string

const (
	FuelTypePetrol   FuelType = "petrol"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeHybrid   FuelType = "hybrid"
)

// TransmissionType enumerates the allowed transmission types.
type TransmissionType string

const (
	TransmissionManual    TransmissionType = "manual"
	TransmissionAutomatic TransmissionType = "automatic"
)

// Vehicle represents a single rental vehicle with all operational, legal, and
// financial data.
type Vehicle struct {
	ID                    uuid.UUID          `json:"id"`
	CarGroupID            uuid.UUID          `json:"car_group_id"`
	BranchID              uuid.UUID          `json:"branch_id"`
	VIN                   string             `json:"vin"`
	LicencePlate          string             `json:"licence_plate"`
	Brand                 string             `json:"brand"`
	Model                 string             `json:"model"`
	Year                  int                `json:"year"`
	Colour                *string            `json:"colour,omitempty"`
	FuelType              FuelType           `json:"fuel_type"`
	TransmissionType      TransmissionType   `json:"transmission_type"`
	CurrentMileage        int                `json:"current_mileage"`
	Status                VehicleStatus      `json:"status"`
	Designation           VehicleDesignation `json:"designation"`
	AcquisitionDate       time.Time          `json:"acquisition_date"`
	OwnershipType         string             `json:"ownership_type"`
	LeaseDetails          *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate   *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time        `json:"registration_expiry_date,omitempty"`
	LastInspectionDate    *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                 *string            `json:"notes,omitempty"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
	DeletedAt             *time.Time         `json:"deleted_at,omitempty"`
	CreatedBy             string             `json:"created_by"`
	UpdatedBy             string             `json:"updated_by"`
	Deleted               bool               `json:"deleted"`
	// ExpiryWarning is true when insurance or registration expires within 30 days.
	ExpiryWarning bool `json:"expiry_warning"`
}

// CreateVehicleInput carries the data required to create a new vehicle record.
type CreateVehicleInput struct {
	CarGroupID              uuid.UUID          `json:"car_group_id"`
	BranchID                uuid.UUID          `json:"branch_id"`
	VIN                     string             `json:"vin"`
	LicencePlate            string             `json:"licence_plate"`
	Brand                   string             `json:"brand"`
	Model                   string             `json:"model"`
	Year                    int                `json:"year"`
	Colour                  *string            `json:"colour,omitempty"`
	FuelType                FuelType           `json:"fuel_type"`
	TransmissionType        TransmissionType   `json:"transmission_type"`
	CurrentMileage          int                `json:"current_mileage"`
	Status                  VehicleStatus      `json:"status"`
	Designation             VehicleDesignation `json:"designation"`
	AcquisitionDate         time.Time          `json:"acquisition_date"`
	OwnershipType           string             `json:"ownership_type"`
	LeaseDetails            *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber   *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate     *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate  *time.Time         `json:"registration_expiry_date,omitempty"`
	LastInspectionDate      *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate   *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                   *string            `json:"notes,omitempty"`
	CreatedBy               string             `json:"-"`
}

// UpdateVehicleInput carries the fields that may be changed on an existing vehicle.
// Status is intentionally excluded; it is managed via dedicated workflow actions.
type UpdateVehicleInput struct {
	CarGroupID              *uuid.UUID         `json:"car_group_id,omitempty"`
	BranchID                *uuid.UUID         `json:"branch_id,omitempty"`
	VIN                     *string            `json:"vin,omitempty"`
	LicencePlate            *string            `json:"licence_plate,omitempty"`
	Brand                   *string            `json:"brand,omitempty"`
	Model                   *string            `json:"model,omitempty"`
	Year                    *int               `json:"year,omitempty"`
	Colour                  *string            `json:"colour,omitempty"`
	FuelType                *FuelType          `json:"fuel_type,omitempty"`
	TransmissionType        *TransmissionType  `json:"transmission_type,omitempty"`
	CurrentMileage          *int               `json:"current_mileage,omitempty"`
	Designation             *VehicleDesignation `json:"designation,omitempty"`
	AcquisitionDate         *time.Time         `json:"acquisition_date,omitempty"`
	OwnershipType           *string            `json:"ownership_type,omitempty"`
	LeaseDetails            *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber   *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate     *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate  *time.Time         `json:"registration_expiry_date,omitempty"`
	LastInspectionDate      *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate   *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                   *string            `json:"notes,omitempty"`
	UpdatedBy               string             `json:"-"`
}

// UpdateDesignationInput carries the new designation value.
type UpdateDesignationInput struct {
	Designation VehicleDesignation `json:"designation"`
	UpdatedBy   string             `json:"-"`
}

// ListVehiclesFilter defines optional filters for listing vehicles.
type ListVehiclesFilter struct {
	CarGroupID       *uuid.UUID
	BranchID         *uuid.UUID
	Status           *VehicleStatus
	Designation      *VehicleDesignation
	FuelType         *FuelType
	TransmissionType *TransmissionType
	// ExpiryWarning filters to vehicles with insurance or registration expiry ≤ 30 days.
	ExpiryWarning *bool
	Page          int
	PageSize      int
}

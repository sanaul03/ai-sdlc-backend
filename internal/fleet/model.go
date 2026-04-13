package fleet

import (
	"time"

	"github.com/google/uuid"
)

// CarGroup represents a vehicle classification group (e.g., Economy Sedan, Luxury SUV).
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

// VehicleStatus represents the operational status of a vehicle.
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

// VehicleDesignation represents the availability designation of a vehicle.
type VehicleDesignation string

const (
	VehicleDesignationRentalOnly VehicleDesignation = "rental_only"
	VehicleDesignationSalesOnly  VehicleDesignation = "sales_only"
	VehicleDesignationShared     VehicleDesignation = "shared"
)

// FuelType represents the fuel type of a vehicle.
type FuelType string

const (
	FuelTypePetrol   FuelType = "petrol"
	FuelTypeDiesel   FuelType = "diesel"
	FuelTypeElectric FuelType = "electric"
	FuelTypeHybrid   FuelType = "hybrid"
)

// TransmissionType represents the transmission type of a vehicle.
type TransmissionType string

const (
	TransmissionTypeManual    TransmissionType = "manual"
	TransmissionTypeAutomatic TransmissionType = "automatic"
)

// OwnershipType represents the ownership model of a vehicle.
type OwnershipType string

const (
	OwnershipTypeOwned  OwnershipType = "owned"
	OwnershipTypeLeased OwnershipType = "leased"
)

// Vehicle stores comprehensive operational, legal, and financial data for each rental vehicle.
type Vehicle struct {
	ID                     uuid.UUID          `json:"id"`
	CarGroupID             uuid.UUID          `json:"car_group_id"`
	BranchID               uuid.UUID          `json:"branch_id"`
	VIN                    string             `json:"vin"`
	LicencePlate           string             `json:"licence_plate"`
	Brand                  string             `json:"brand"`
	Model                  string             `json:"model"`
	Year                   int                `json:"year"`
	Colour                 *string            `json:"colour,omitempty"`
	FuelType               FuelType           `json:"fuel_type"`
	TransmissionType       TransmissionType   `json:"transmission_type"`
	CurrentMileage         int                `json:"current_mileage"`
	Status                 VehicleStatus      `json:"status"`
	Designation            VehicleDesignation `json:"designation"`
	AcquisitionDate        time.Time          `json:"acquisition_date"`
	OwnershipType          OwnershipType      `json:"ownership_type"`
	LeaseDetails           *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time         `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                  *string            `json:"notes,omitempty"`
	CreatedAt              time.Time          `json:"created_at"`
	UpdatedAt              time.Time          `json:"updated_at"`
	DeletedAt              *time.Time         `json:"deleted_at,omitempty"`
	CreatedBy              string             `json:"created_by"`
	UpdatedBy              string             `json:"updated_by"`
	Deleted                bool               `json:"deleted"`
}

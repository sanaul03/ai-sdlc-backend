package fleet

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// vinRegex is the standard 17-character VIN pattern (excludes I, O, Q).
var vinRegex = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// VehicleService implements business logic for vehicle operations.
type VehicleService struct {
	repo VehicleRepository
}

// NewVehicleService constructs a VehicleService backed by the given repository.
func NewVehicleService(repo VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

// Create validates the input and creates a new vehicle record.
func (s *VehicleService) Create(ctx context.Context, input CreateVehicleInput) (*Vehicle, error) {
	if err := validateVehicleCreate(input); err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, input)
}

// List retrieves vehicles matching the filter.
func (s *VehicleService) List(ctx context.Context, filter ListVehiclesFilter) ([]*Vehicle, int, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	return s.repo.List(ctx, filter)
}

// GetByID retrieves a single vehicle by its identifier.
func (s *VehicleService) GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}

// Update validates the input and applies changes to an existing vehicle.
func (s *VehicleService) Update(ctx context.Context, id uuid.UUID, input UpdateVehicleInput) (*Vehicle, error) {
	if err := validateVehicleUpdate(input); err != nil {
		return nil, err
	}
	return s.repo.Update(ctx, id, input)
}

// UpdateDesignation changes the designation of a vehicle.
func (s *VehicleService) UpdateDesignation(ctx context.Context, id uuid.UUID, input UpdateDesignationInput) (*Vehicle, error) {
	if err := validateDesignation(input.Designation); err != nil {
		return nil, err
	}
	return s.repo.UpdateDesignation(ctx, id, input)
}

// Delete soft-deletes a vehicle.
func (s *VehicleService) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return s.repo.Delete(ctx, id, deletedBy)
}

// validateVehicleCreate checks required fields when creating a vehicle.
func validateVehicleCreate(input CreateVehicleInput) error {
	if input.CarGroupID == uuid.Nil {
		return fmt.Errorf("%w: car_group_id is required", ErrValidation)
	}
	if input.BranchID == uuid.Nil {
		return fmt.Errorf("%w: branch_id is required", ErrValidation)
	}
	if !vinRegex.MatchString(input.VIN) {
		return fmt.Errorf("%w: vin must be a valid 17-character VIN", ErrValidation)
	}
	if strings.TrimSpace(input.LicencePlate) == "" {
		return fmt.Errorf("%w: licence_plate is required", ErrValidation)
	}
	if len(strings.TrimSpace(input.LicencePlate)) > 20 {
		return fmt.Errorf("%w: licence_plate must not exceed 20 characters", ErrValidation)
	}
	if strings.TrimSpace(input.Brand) == "" {
		return fmt.Errorf("%w: brand is required", ErrValidation)
	}
	if strings.TrimSpace(input.Model) == "" {
		return fmt.Errorf("%w: model is required", ErrValidation)
	}
	currentYear := time.Now().Year()
	if input.Year < 1900 || input.Year > currentYear+1 {
		return fmt.Errorf("%w: year must be between 1900 and %d", ErrValidation, currentYear+1)
	}
	if err := validateFuelType(input.FuelType); err != nil {
		return err
	}
	if err := validateTransmissionType(input.TransmissionType); err != nil {
		return err
	}
	if input.CurrentMileage < 0 {
		return fmt.Errorf("%w: current_mileage must be non-negative", ErrValidation)
	}
	if err := validateVehicleStatus(input.Status); err != nil {
		return err
	}
	if err := validateDesignation(input.Designation); err != nil {
		return err
	}
	if input.AcquisitionDate.IsZero() {
		return fmt.Errorf("%w: acquisition_date is required", ErrValidation)
	}
	if input.AcquisitionDate.After(time.Now()) {
		return fmt.Errorf("%w: acquisition_date must not be a future date", ErrValidation)
	}
	if strings.TrimSpace(input.OwnershipType) == "" {
		return fmt.Errorf("%w: ownership_type is required", ErrValidation)
	}
	if strings.TrimSpace(input.CreatedBy) == "" {
		return fmt.Errorf("%w: created_by is required", ErrValidation)
	}

	// When status is 'available', mandatory compliance fields must be present.
	if input.Status == VehicleStatusAvailable {
		if err := validateAvailabilityRequirements(
			input.InsurancePolicyNumber,
			input.InsuranceExpiryDate,
			input.RegistrationExpiryDate,
		); err != nil {
			return err
		}
	}
	return nil
}

// validateVehicleUpdate validates fields supplied in an update request.
func validateVehicleUpdate(input UpdateVehicleInput) error {
	if input.VIN != nil && !vinRegex.MatchString(*input.VIN) {
		return fmt.Errorf("%w: vin must be a valid 17-character VIN", ErrValidation)
	}
	if input.LicencePlate != nil {
		if strings.TrimSpace(*input.LicencePlate) == "" {
			return fmt.Errorf("%w: licence_plate must not be empty", ErrValidation)
		}
		if len(strings.TrimSpace(*input.LicencePlate)) > 20 {
			return fmt.Errorf("%w: licence_plate must not exceed 20 characters", ErrValidation)
		}
	}
	if input.Year != nil {
		currentYear := time.Now().Year()
		if *input.Year < 1900 || *input.Year > currentYear+1 {
			return fmt.Errorf("%w: year must be between 1900 and %d", ErrValidation, currentYear+1)
		}
	}
	if input.FuelType != nil {
		if err := validateFuelType(*input.FuelType); err != nil {
			return err
		}
	}
	if input.TransmissionType != nil {
		if err := validateTransmissionType(*input.TransmissionType); err != nil {
			return err
		}
	}
	if input.CurrentMileage != nil && *input.CurrentMileage < 0 {
		return fmt.Errorf("%w: current_mileage must be non-negative", ErrValidation)
	}
	if input.Designation != nil {
		if err := validateDesignation(*input.Designation); err != nil {
			return err
		}
	}
	if input.AcquisitionDate != nil && input.AcquisitionDate.After(time.Now()) {
		return fmt.Errorf("%w: acquisition_date must not be a future date", ErrValidation)
	}
	return nil
}

func validateFuelType(ft FuelType) error {
	switch ft {
	case FuelTypePetrol, FuelTypeDiesel, FuelTypeElectric, FuelTypeHybrid:
		return nil
	default:
		return fmt.Errorf("%w: fuel_type must be one of: petrol, diesel, electric, hybrid", ErrValidation)
	}
}

func validateTransmissionType(tt TransmissionType) error {
	switch tt {
	case TransmissionManual, TransmissionAutomatic:
		return nil
	default:
		return fmt.Errorf("%w: transmission_type must be one of: manual, automatic", ErrValidation)
	}
}

func validateVehicleStatus(s VehicleStatus) error {
	switch s {
	case VehicleStatusAvailable, VehicleStatusOnRent, VehicleStatusNeedsCleaning,
		VehicleStatusNeedsInspection, VehicleStatusUnderMaintenance,
		VehicleStatusUnavailable, VehicleStatusDecommissioned:
		return nil
	default:
		return fmt.Errorf("%w: status must be one of: available, on_rent, needs_cleaning, needs_inspection, under_maintenance, unavailable, decommissioned", ErrValidation)
	}
}

func validateDesignation(d VehicleDesignation) error {
	switch d {
	case DesignationRentalOnly, DesignationSalesOnly, DesignationShared:
		return nil
	default:
		return fmt.Errorf("%w: designation must be one of: rental_only, sales_only, shared", ErrValidation)
	}
}

func validateAvailabilityRequirements(policyNumber *string, insuranceExpiry, registrationExpiry *time.Time) error {
	if policyNumber == nil || strings.TrimSpace(*policyNumber) == "" {
		return fmt.Errorf("%w: insurance_policy_number is required when status is available", ErrValidation)
	}
	if insuranceExpiry == nil {
		return fmt.Errorf("%w: insurance_expiry_date is required when status is available", ErrValidation)
	}
	if registrationExpiry == nil {
		return fmt.Errorf("%w: registration_expiry_date is required when status is available", ErrValidation)
	}
	return nil
}

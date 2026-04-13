package vehicle

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// vinRegex validates a standard 17-character VIN (excludes I, O, Q).
var vinRegex = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// Service provides business logic for vehicles.
type Service interface {
	Create(ctx context.Context, req CreateRequest, createdBy string) (Vehicle, error)
	List(ctx context.Context, filter ListFilter) (Page, error)
	GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest, updatedBy string) (Vehicle, error)
	UpdateDesignation(ctx context.Context, id uuid.UUID, req DesignationUpdateRequest, updatedBy string) (Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

type service struct {
	repo Repository
}

// NewService returns a Service backed by the provided Repository.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateRequest, createdBy string) (Vehicle, error) {
	if err := validateCreateRequest(req); err != nil {
		return Vehicle{}, err
	}

	now := time.Now().UTC()
	v := Vehicle{
		ID:                     uuid.New(),
		CarGroupID:             req.CarGroupID,
		BranchID:               req.BranchID,
		VIN:                    strings.ToUpper(strings.TrimSpace(req.VIN)),
		LicencePlate:           strings.TrimSpace(req.LicencePlate),
		Brand:                  strings.TrimSpace(req.Brand),
		Model:                  strings.TrimSpace(req.Model),
		Year:                   req.Year,
		Colour:                 req.Colour,
		FuelType:               req.FuelType,
		TransmissionType:       req.TransmissionType,
		CurrentMileage:         req.CurrentMileage,
		Status:                 req.Status,
		Designation:            req.Designation,
		AcquisitionDate:        req.AcquisitionDate,
		OwnershipType:          req.OwnershipType,
		LeaseDetails:           req.LeaseDetails,
		InsurancePolicyNumber:  req.InsurancePolicyNumber,
		InsuranceExpiryDate:    req.InsuranceExpiryDate,
		RegistrationExpiryDate: req.RegistrationExpiryDate,
		LastInspectionDate:     req.LastInspectionDate,
		NextInspectionDueDate:  req.NextInspectionDueDate,
		Notes:                  req.Notes,
		CreatedAt:              now,
		UpdatedAt:              now,
		CreatedBy:              createdBy,
		UpdatedBy:              createdBy,
		Deleted:                false,
	}

	if err := validateAvailabilityFields(v); err != nil {
		return Vehicle{}, err
	}

	return s.repo.Create(ctx, v)
}

func (s *service) List(ctx context.Context, filter ListFilter) (Page, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest, updatedBy string) (Vehicle, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Vehicle{}, err
	}

	applyUpdateRequest(&existing, req)
	existing.UpdatedAt = time.Now().UTC()
	existing.UpdatedBy = updatedBy

	if err := validateUpdateFields(existing); err != nil {
		return Vehicle{}, err
	}
	if err := validateAvailabilityFields(existing); err != nil {
		return Vehicle{}, err
	}

	return s.repo.Update(ctx, existing)
}

func (s *service) UpdateDesignation(ctx context.Context, id uuid.UUID, req DesignationUpdateRequest, updatedBy string) (Vehicle, error) {
	if _, ok := AllowedDesignations[req.Designation]; !ok {
		return Vehicle{}, fmt.Errorf("%w: designation must be one of rental_only, sales_only, shared", ErrInvalidInput)
	}
	return s.repo.UpdateDesignation(ctx, id, req.Designation, updatedBy)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return s.repo.SoftDelete(ctx, id, deletedBy)
}

// --- validation helpers ---

func validateCreateRequest(req CreateRequest) error {
	if req.CarGroupID == uuid.Nil {
		return fmt.Errorf("%w: car_group_id is required", ErrInvalidInput)
	}
	if req.BranchID == uuid.Nil {
		return fmt.Errorf("%w: branch_id is required", ErrInvalidInput)
	}
	vin := strings.ToUpper(strings.TrimSpace(req.VIN))
	if !vinRegex.MatchString(vin) {
		return fmt.Errorf("%w: vin must be a valid 17-character VIN", ErrInvalidInput)
	}
	if strings.TrimSpace(req.LicencePlate) == "" {
		return fmt.Errorf("%w: licence_plate is required", ErrInvalidInput)
	}
	if utf8.RuneCountInString(strings.TrimSpace(req.LicencePlate)) > 20 {
		return fmt.Errorf("%w: licence_plate must be 20 characters or fewer", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Brand) == "" {
		return fmt.Errorf("%w: brand is required", ErrInvalidInput)
	}
	if strings.TrimSpace(req.Model) == "" {
		return fmt.Errorf("%w: model is required", ErrInvalidInput)
	}
	currentYear := time.Now().Year()
	if req.Year < 1900 || req.Year > currentYear+1 {
		return fmt.Errorf("%w: year must be between 1900 and %d", ErrInvalidInput, currentYear+1)
	}
	if _, ok := AllowedFuelTypes[req.FuelType]; !ok {
		return fmt.Errorf("%w: fuel_type must be one of petrol, diesel, electric, hybrid", ErrInvalidInput)
	}
	if _, ok := AllowedTransmissionTypes[req.TransmissionType]; !ok {
		return fmt.Errorf("%w: transmission_type must be one of manual, automatic", ErrInvalidInput)
	}
	if req.CurrentMileage < 0 {
		return fmt.Errorf("%w: current_mileage must be non-negative", ErrInvalidInput)
	}
	if _, ok := AllowedStatuses[req.Status]; !ok {
		return fmt.Errorf("%w: status is invalid", ErrInvalidInput)
	}
	if _, ok := AllowedDesignations[req.Designation]; !ok {
		return fmt.Errorf("%w: designation must be one of rental_only, sales_only, shared", ErrInvalidInput)
	}
	if req.AcquisitionDate.IsZero() {
		return fmt.Errorf("%w: acquisition_date is required", ErrInvalidInput)
	}
	if req.AcquisitionDate.After(time.Now().UTC()) {
		return fmt.Errorf("%w: acquisition_date must not be a future date", ErrInvalidInput)
	}
	if strings.TrimSpace(req.OwnershipType) == "" {
		return fmt.Errorf("%w: ownership_type is required", ErrInvalidInput)
	}
	return nil
}

func validateUpdateFields(v Vehicle) error {
	vin := strings.ToUpper(strings.TrimSpace(v.VIN))
	if !vinRegex.MatchString(vin) {
		return fmt.Errorf("%w: vin must be a valid 17-character VIN", ErrInvalidInput)
	}
	if strings.TrimSpace(v.LicencePlate) == "" {
		return fmt.Errorf("%w: licence_plate is required", ErrInvalidInput)
	}
	if _, ok := AllowedFuelTypes[v.FuelType]; !ok {
		return fmt.Errorf("%w: fuel_type must be one of petrol, diesel, electric, hybrid", ErrInvalidInput)
	}
	if _, ok := AllowedTransmissionTypes[v.TransmissionType]; !ok {
		return fmt.Errorf("%w: transmission_type must be one of manual, automatic", ErrInvalidInput)
	}
	if _, ok := AllowedDesignations[v.Designation]; !ok {
		return fmt.Errorf("%w: designation must be one of rental_only, sales_only, shared", ErrInvalidInput)
	}
	if v.CurrentMileage < 0 {
		return fmt.Errorf("%w: current_mileage must be non-negative", ErrInvalidInput)
	}
	return nil
}

// validateAvailabilityFields enforces the rule that a vehicle can only be set to
// `available` when all mandatory fields are populated.
func validateAvailabilityFields(v Vehicle) error {
	if v.Status != "available" {
		return nil
	}
	if v.VIN == "" {
		return fmt.Errorf("%w: vin is required before setting status to available", ErrInvalidInput)
	}
	if v.LicencePlate == "" {
		return fmt.Errorf("%w: licence_plate is required before setting status to available", ErrInvalidInput)
	}
	if v.InsurancePolicyNumber == nil || *v.InsurancePolicyNumber == "" {
		return fmt.Errorf("%w: insurance_policy_number is required before setting status to available", ErrInvalidInput)
	}
	if v.InsuranceExpiryDate == nil {
		return fmt.Errorf("%w: insurance_expiry_date is required before setting status to available", ErrInvalidInput)
	}
	if v.RegistrationExpiryDate == nil {
		return fmt.Errorf("%w: registration_expiry_date is required before setting status to available", ErrInvalidInput)
	}
	return nil
}

// applyUpdateRequest merges non-nil fields from UpdateRequest onto an existing Vehicle.
func applyUpdateRequest(v *Vehicle, req UpdateRequest) {
	if req.CarGroupID != nil {
		v.CarGroupID = *req.CarGroupID
	}
	if req.BranchID != nil {
		v.BranchID = *req.BranchID
	}
	if req.VIN != nil {
		v.VIN = strings.ToUpper(strings.TrimSpace(*req.VIN))
	}
	if req.LicencePlate != nil {
		v.LicencePlate = strings.TrimSpace(*req.LicencePlate)
	}
	if req.Brand != nil {
		v.Brand = strings.TrimSpace(*req.Brand)
	}
	if req.Model != nil {
		v.Model = strings.TrimSpace(*req.Model)
	}
	if req.Year != nil {
		v.Year = *req.Year
	}
	if req.Colour != nil {
		v.Colour = req.Colour
	}
	if req.FuelType != nil {
		v.FuelType = *req.FuelType
	}
	if req.TransmissionType != nil {
		v.TransmissionType = *req.TransmissionType
	}
	if req.CurrentMileage != nil {
		v.CurrentMileage = *req.CurrentMileage
	}
	if req.Designation != nil {
		v.Designation = *req.Designation
	}
	if req.AcquisitionDate != nil {
		v.AcquisitionDate = *req.AcquisitionDate
	}
	if req.OwnershipType != nil {
		v.OwnershipType = *req.OwnershipType
	}
	if req.LeaseDetails != nil {
		v.LeaseDetails = req.LeaseDetails
	}
	if req.InsurancePolicyNumber != nil {
		v.InsurancePolicyNumber = req.InsurancePolicyNumber
	}
	if req.InsuranceExpiryDate != nil {
		v.InsuranceExpiryDate = req.InsuranceExpiryDate
	}
	if req.RegistrationExpiryDate != nil {
		v.RegistrationExpiryDate = req.RegistrationExpiryDate
	}
	if req.LastInspectionDate != nil {
		v.LastInspectionDate = req.LastInspectionDate
	}
	if req.NextInspectionDueDate != nil {
		v.NextInspectionDueDate = req.NextInspectionDueDate
	}
	if req.Notes != nil {
		v.Notes = req.Notes
	}
}

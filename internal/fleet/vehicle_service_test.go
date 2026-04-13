package fleet_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// --- mock repository ---

type mockVehicleRepo struct {
	createFn            func(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error)
	listFn              func(ctx context.Context, filter fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error)
	getByIDFn           func(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error)
	updateFn            func(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error)
	updateDesignationFn func(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error)
	deleteFn            func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockVehicleRepo) Create(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
	return m.createFn(ctx, input)
}
func (m *mockVehicleRepo) List(ctx context.Context, filter fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error) {
	return m.listFn(ctx, filter)
}
func (m *mockVehicleRepo) GetByID(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockVehicleRepo) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error) {
	return m.updateFn(ctx, id, input)
}
func (m *mockVehicleRepo) UpdateDesignation(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error) {
	return m.updateDesignationFn(ctx, id, input)
}
func (m *mockVehicleRepo) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.deleteFn(ctx, id, deletedBy)
}

// helpers

func validCreateVehicleInput() fleet.CreateVehicleInput {
	acq := time.Now().AddDate(-1, 0, 0)
	return fleet.CreateVehicleInput{
		CarGroupID:       uuid.New(),
		BranchID:         uuid.New(),
		VIN:              "1HGBH41JXMN109186",
		LicencePlate:     "ABC-1234",
		Brand:            "Toyota",
		Model:            "Corolla",
		Year:             2022,
		FuelType:         fleet.FuelTypePetrol,
		TransmissionType: fleet.TransmissionAutomatic,
		CurrentMileage:   0,
		Status:           fleet.VehicleStatusUnavailable,
		Designation:      fleet.DesignationRentalOnly,
		AcquisitionDate:  acq,
		OwnershipType:    "owned",
		CreatedBy:        "admin",
	}
}

// --- tests ---

func TestVehicleService_Create_Success(t *testing.T) {
	id := uuid.New()
	repo := &mockVehicleRepo{
		createFn: func(_ context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
			return &fleet.Vehicle{ID: id, VIN: input.VIN}, nil
		},
	}
	svc := fleet.NewVehicleService(repo)
	got, err := svc.Create(context.Background(), validCreateVehicleInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("expected id %v, got %v", id, got.ID)
	}
}

func TestVehicleService_Create_InvalidVIN(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.VIN = "INVALID"
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for invalid VIN, got %v", err)
	}
}

func TestVehicleService_Create_FutureAcquisitionDate(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.AcquisitionDate = time.Now().AddDate(1, 0, 0)
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for future acquisition date, got %v", err)
	}
}

func TestVehicleService_Create_InvalidFuelType(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.FuelType = "gasoline"
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for invalid fuel type, got %v", err)
	}
}

func TestVehicleService_Create_InvalidTransmission(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.TransmissionType = "cvt"
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for invalid transmission type, got %v", err)
	}
}

func TestVehicleService_Create_NegativeMileage(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.CurrentMileage = -1
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for negative mileage, got %v", err)
	}
}

func TestVehicleService_Create_StatusAvailable_MissingInsurance(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.Status = fleet.VehicleStatusAvailable
	// InsurancePolicyNumber, InsuranceExpiryDate, RegistrationExpiryDate are nil
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation when available without insurance, got %v", err)
	}
}

func TestVehicleService_Create_StatusAvailable_AllRequiredFields(t *testing.T) {
	policy := "POL-001"
	exp := time.Now().AddDate(1, 0, 0)
	regExp := time.Now().AddDate(1, 0, 0)
	id := uuid.New()
	repo := &mockVehicleRepo{
		createFn: func(_ context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
			return &fleet.Vehicle{ID: id, Status: fleet.VehicleStatusAvailable}, nil
		},
	}
	svc := fleet.NewVehicleService(repo)
	input := validCreateVehicleInput()
	input.Status = fleet.VehicleStatusAvailable
	input.InsurancePolicyNumber = &policy
	input.InsuranceExpiryDate = &exp
	input.RegistrationExpiryDate = &regExp
	got, err := svc.Create(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Status != fleet.VehicleStatusAvailable {
		t.Errorf("expected status available, got %v", got.Status)
	}
}

func TestVehicleService_Create_InvalidDesignation(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.Designation = "private"
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for invalid designation, got %v", err)
	}
}

func TestVehicleService_Create_YearOutOfRange(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	input := validCreateVehicleInput()
	input.Year = 1800
	_, err := svc.Create(context.Background(), input)
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation for year out of range, got %v", err)
	}
}

func TestVehicleService_UpdateDesignation_Invalid(t *testing.T) {
	svc := fleet.NewVehicleService(&mockVehicleRepo{})
	_, err := svc.UpdateDesignation(context.Background(), uuid.New(), fleet.UpdateDesignationInput{
		Designation: "unknown",
	})
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestVehicleService_List_DefaultPagination(t *testing.T) {
	var capturedFilter fleet.ListVehiclesFilter
	repo := &mockVehicleRepo{
		listFn: func(_ context.Context, f fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error) {
			capturedFilter = f
			return nil, 0, nil
		},
	}
	svc := fleet.NewVehicleService(repo)
	_, _, err := svc.List(context.Background(), fleet.ListVehiclesFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedFilter.Page != 1 {
		t.Errorf("expected default page 1, got %d", capturedFilter.Page)
	}
	if capturedFilter.PageSize != 20 {
		t.Errorf("expected default page size 20, got %d", capturedFilter.PageSize)
	}
}

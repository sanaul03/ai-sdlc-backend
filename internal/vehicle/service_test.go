package vehicle_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/vehicle"
)

// --- mock repository ---

type mockRepo struct {
	createFn            func(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error)
	listFn              func(ctx context.Context, filter vehicle.ListFilter) (vehicle.Page, error)
	getByIDFn           func(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error)
	updateFn            func(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error)
	updateDesignationFn func(ctx context.Context, id uuid.UUID, designation, updatedBy string) (vehicle.Vehicle, error)
	softDeleteFn        func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockRepo) Create(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	return m.createFn(ctx, v)
}
func (m *mockRepo) List(ctx context.Context, f vehicle.ListFilter) (vehicle.Page, error) {
	return m.listFn(ctx, f)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockRepo) Update(ctx context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
	return m.updateFn(ctx, v)
}
func (m *mockRepo) UpdateDesignation(ctx context.Context, id uuid.UUID, designation, updatedBy string) (vehicle.Vehicle, error) {
	return m.updateDesignationFn(ctx, id, designation, updatedBy)
}
func (m *mockRepo) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.softDeleteFn(ctx, id, deletedBy)
}

// --- helpers ---

func newValidCreateRequest() vehicle.CreateRequest {
	acquiredAt := time.Now().UTC().Add(-24 * time.Hour)
	return vehicle.CreateRequest{
		CarGroupID:       uuid.New(),
		BranchID:         uuid.New(),
		VIN:              "1HGBH41JXMN109186",
		LicencePlate:     "ABC-1234",
		Brand:            "Toyota",
		Model:            "Corolla",
		Year:             2020,
		FuelType:         "petrol",
		TransmissionType: "automatic",
		CurrentMileage:   0,
		Status:           "unavailable",
		Designation:      "rental_only",
		AcquisitionDate:  acquiredAt,
		OwnershipType:    "owned",
	}
}

func newSampleVehicle() vehicle.Vehicle {
	acquiredAt := time.Now().UTC().Add(-24 * time.Hour)
	return vehicle.Vehicle{
		ID:               uuid.New(),
		CarGroupID:       uuid.New(),
		BranchID:         uuid.New(),
		VIN:              "1HGBH41JXMN109186",
		LicencePlate:     "ABC-1234",
		Brand:            "Toyota",
		Model:            "Corolla",
		Year:             2020,
		FuelType:         "petrol",
		TransmissionType: "automatic",
		CurrentMileage:   0,
		Status:           "unavailable",
		Designation:      "rental_only",
		AcquisitionDate:  acquiredAt,
		OwnershipType:    "owned",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		CreatedBy:        "user1",
		UpdatedBy:        "user1",
	}
}

// --- tests ---

func TestService_Create_Success(t *testing.T) {
	repo := &mockRepo{
		createFn: func(_ context.Context, v vehicle.Vehicle) (vehicle.Vehicle, error) {
			return v, nil
		},
	}
	svc := vehicle.NewService(repo)

	req := newValidCreateRequest()
	got, err := svc.Create(context.Background(), req, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.VIN != req.VIN {
		t.Errorf("expected VIN %q, got %q", req.VIN, got.VIN)
	}
}

func TestService_Create_InvalidVIN(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	req := newValidCreateRequest()
	req.VIN = "INVALID"

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Create_InvalidFuelType(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	req := newValidCreateRequest()
	req.FuelType = "coal"

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Create_InvalidTransmission(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	req := newValidCreateRequest()
	req.TransmissionType = "semi"

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Create_FutureAcquisitionDate(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	req := newValidCreateRequest()
	req.AcquisitionDate = time.Now().UTC().Add(48 * time.Hour)

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for future acquisition date, got %v", err)
	}
}

func TestService_Create_StatusAvailableMissingFields(t *testing.T) {
	repo := &mockRepo{}
	svc := vehicle.NewService(repo)
	req := newValidCreateRequest()
	req.Status = "available"
	// InsurancePolicyNumber is nil → should fail

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput for missing insurance fields, got %v", err)
	}
}

func TestService_Create_NegativeMileage(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	req := newValidCreateRequest()
	req.CurrentMileage = -1

	_, err := svc.Create(context.Background(), req, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (vehicle.Vehicle, error) {
			return vehicle.Vehicle{}, vehicle.ErrNotFound
		},
	}
	svc := vehicle.NewService(repo)

	_, err := svc.GetByID(context.Background(), uuid.New())
	if !errors.Is(err, vehicle.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_UpdateDesignation_Invalid(t *testing.T) {
	svc := vehicle.NewService(&mockRepo{})
	_, err := svc.UpdateDesignation(context.Background(), uuid.New(),
		vehicle.DesignationUpdateRequest{Designation: "unknown"}, "user1")
	if !errors.Is(err, vehicle.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_UpdateDesignation_Success(t *testing.T) {
	sample := newSampleVehicle()
	sample.Designation = "shared"
	repo := &mockRepo{
		updateDesignationFn: func(_ context.Context, _ uuid.UUID, d, _ string) (vehicle.Vehicle, error) {
			sample.Designation = d
			return sample, nil
		},
	}
	svc := vehicle.NewService(repo)

	got, err := svc.UpdateDesignation(context.Background(), sample.ID,
		vehicle.DesignationUpdateRequest{Designation: "sales_only"}, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Designation != "sales_only" {
		t.Errorf("expected designation %q, got %q", "sales_only", got.Designation)
	}
}

func TestService_Delete_Success(t *testing.T) {
	repo := &mockRepo{
		softDeleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	svc := vehicle.NewService(repo)
	if err := svc.Delete(context.Background(), uuid.New(), "user1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	repo := &mockRepo{
		softDeleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return vehicle.ErrNotFound
		},
	}
	svc := vehicle.NewService(repo)
	err := svc.Delete(context.Background(), uuid.New(), "user1")
	if !errors.Is(err, vehicle.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_List(t *testing.T) {
	expected := vehicle.Page{Items: []vehicle.Vehicle{newSampleVehicle()}, Total: 1, Page: 1, PageSize: 20}
	repo := &mockRepo{
		listFn: func(_ context.Context, _ vehicle.ListFilter) (vehicle.Page, error) {
			return expected, nil
		},
	}
	svc := vehicle.NewService(repo)
	got, err := svc.List(context.Background(), vehicle.ListFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Total != expected.Total {
		t.Errorf("expected total %d, got %d", expected.Total, got.Total)
	}
}

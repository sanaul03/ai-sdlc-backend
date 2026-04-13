package fleet_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// --- mock repository ---

type mockCarGroupRepo struct {
	createFn            func(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error)
	listFn              func(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error)
	getByIDFn           func(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error)
	updateFn            func(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error)
	deleteFn            func(ctx context.Context, id uuid.UUID, deletedBy string) error
	hasActiveVehiclesFn func(ctx context.Context, carGroupID uuid.UUID) (bool, error)
}

func (m *mockCarGroupRepo) Create(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error) {
	return m.createFn(ctx, input)
}
func (m *mockCarGroupRepo) List(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error) {
	return m.listFn(ctx, filter)
}
func (m *mockCarGroupRepo) GetByID(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCarGroupRepo) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error) {
	return m.updateFn(ctx, id, input)
}
func (m *mockCarGroupRepo) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.deleteFn(ctx, id, deletedBy)
}
func (m *mockCarGroupRepo) HasActiveVehicles(ctx context.Context, carGroupID uuid.UUID) (bool, error) {
	return m.hasActiveVehiclesFn(ctx, carGroupID)
}

// --- tests ---

func TestCarGroupService_Create_Success(t *testing.T) {
	id := uuid.New()
	name := "Economy Sedan"
	repo := &mockCarGroupRepo{
		createFn: func(_ context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error) {
			return &fleet.CarGroup{ID: id, Name: input.Name, CreatedBy: input.CreatedBy}, nil
		},
	}
	svc := fleet.NewCarGroupService(repo)

	got, err := svc.Create(context.Background(), fleet.CreateCarGroupInput{
		Name:      name,
		CreatedBy: "admin",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != name {
		t.Errorf("expected name %q, got %q", name, got.Name)
	}
}

func TestCarGroupService_Create_EmptyName(t *testing.T) {
	svc := fleet.NewCarGroupService(&mockCarGroupRepo{})
	_, err := svc.Create(context.Background(), fleet.CreateCarGroupInput{
		Name:      "  ",
		CreatedBy: "admin",
	})
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestCarGroupService_Create_MissingCreatedBy(t *testing.T) {
	svc := fleet.NewCarGroupService(&mockCarGroupRepo{})
	_, err := svc.Create(context.Background(), fleet.CreateCarGroupInput{
		Name:      "Economy",
		CreatedBy: "",
	})
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestCarGroupService_Delete_WithActiveVehicles(t *testing.T) {
	id := uuid.New()
	repo := &mockCarGroupRepo{
		hasActiveVehiclesFn: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
	}
	svc := fleet.NewCarGroupService(repo)
	err := svc.Delete(context.Background(), id, "admin")
	if !errors.Is(err, fleet.ErrConflict) {
		t.Errorf("expected ErrConflict, got %v", err)
	}
}

func TestCarGroupService_Delete_NoActiveVehicles(t *testing.T) {
	id := uuid.New()
	deleted := false
	repo := &mockCarGroupRepo{
		hasActiveVehiclesFn: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return false, nil
		},
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			deleted = true
			return nil
		},
	}
	svc := fleet.NewCarGroupService(repo)
	if err := svc.Delete(context.Background(), id, "admin"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("expected Delete to be called on repository")
	}
}

func TestCarGroupService_Update_EmptyName(t *testing.T) {
	empty := ""
	svc := fleet.NewCarGroupService(&mockCarGroupRepo{})
	_, err := svc.Update(context.Background(), uuid.New(), fleet.UpdateCarGroupInput{
		Name: &empty,
	})
	if !errors.Is(err, fleet.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

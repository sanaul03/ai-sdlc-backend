package cargroup_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/cargroup"
)

// --- mock repository ---

type mockRepo struct {
	createFn          func(ctx context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error)
	listFn            func(ctx context.Context, filter cargroup.ListFilter) ([]cargroup.CarGroup, error)
	getByIDFn         func(ctx context.Context, id uuid.UUID) (cargroup.CarGroup, error)
	updateFn          func(ctx context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error)
	softDeleteFn      func(ctx context.Context, id uuid.UUID, deletedBy string) error
	hasActiveVehicles func(ctx context.Context, id uuid.UUID) (bool, error)
}

func (m *mockRepo) Create(ctx context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error) {
	return m.createFn(ctx, cg)
}
func (m *mockRepo) List(ctx context.Context, filter cargroup.ListFilter) ([]cargroup.CarGroup, error) {
	return m.listFn(ctx, filter)
}
func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (cargroup.CarGroup, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockRepo) Update(ctx context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error) {
	return m.updateFn(ctx, cg)
}
func (m *mockRepo) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.softDeleteFn(ctx, id, deletedBy)
}
func (m *mockRepo) HasActiveVehicles(ctx context.Context, id uuid.UUID) (bool, error) {
	return m.hasActiveVehicles(ctx, id)
}

// --- helpers ---

func newSampleGroup() cargroup.CarGroup {
	desc := "Economy class"
	size := "compact"
	return cargroup.CarGroup{
		ID:           uuid.New(),
		Name:         "Economy Sedan",
		Description:  &desc,
		SizeCategory: &size,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		CreatedBy:    "user1",
		UpdatedBy:    "user1",
		Deleted:      false,
	}
}

// --- tests ---

func TestService_Create_Success(t *testing.T) {
	sample := newSampleGroup()
	repo := &mockRepo{
		createFn: func(_ context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error) {
			return cg, nil
		},
	}
	svc := cargroup.NewService(repo)

	got, err := svc.Create(context.Background(), cargroup.CreateRequest{Name: sample.Name}, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != sample.Name {
		t.Errorf("expected name %q, got %q", sample.Name, got.Name)
	}
}

func TestService_Create_EmptyName(t *testing.T) {
	svc := cargroup.NewService(&mockRepo{})

	_, err := svc.Create(context.Background(), cargroup.CreateRequest{Name: "  "}, "user1")
	if !errors.Is(err, cargroup.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (cargroup.CarGroup, error) {
			return cargroup.CarGroup{}, cargroup.ErrNotFound
		},
	}
	svc := cargroup.NewService(repo)

	_, err := svc.GetByID(context.Background(), uuid.New())
	if !errors.Is(err, cargroup.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestService_Update_EmptyName(t *testing.T) {
	sample := newSampleGroup()
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (cargroup.CarGroup, error) {
			return sample, nil
		},
	}
	svc := cargroup.NewService(repo)

	emptyName := "  "
	_, err := svc.Update(context.Background(), sample.ID, cargroup.UpdateRequest{Name: &emptyName}, "user1")
	if !errors.Is(err, cargroup.ErrInvalidInput) {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestService_Update_Success(t *testing.T) {
	sample := newSampleGroup()
	newName := "Luxury SUV"
	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (cargroup.CarGroup, error) {
			return sample, nil
		},
		updateFn: func(_ context.Context, cg cargroup.CarGroup) (cargroup.CarGroup, error) {
			return cg, nil
		},
	}
	svc := cargroup.NewService(repo)

	got, err := svc.Update(context.Background(), sample.ID, cargroup.UpdateRequest{Name: &newName}, "user1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != newName {
		t.Errorf("expected name %q, got %q", newName, got.Name)
	}
}

func TestService_Delete_HasVehicles(t *testing.T) {
	repo := &mockRepo{
		hasActiveVehicles: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return true, nil
		},
	}
	svc := cargroup.NewService(repo)

	err := svc.Delete(context.Background(), uuid.New(), "user1")
	if !errors.Is(err, cargroup.ErrHasVehicles) {
		t.Fatalf("expected ErrHasVehicles, got %v", err)
	}
}

func TestService_Delete_Success(t *testing.T) {
	repo := &mockRepo{
		hasActiveVehicles: func(_ context.Context, _ uuid.UUID) (bool, error) {
			return false, nil
		},
		softDeleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	svc := cargroup.NewService(repo)

	if err := svc.Delete(context.Background(), uuid.New(), "user1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_List(t *testing.T) {
	expected := []cargroup.CarGroup{newSampleGroup()}
	repo := &mockRepo{
		listFn: func(_ context.Context, _ cargroup.ListFilter) ([]cargroup.CarGroup, error) {
			return expected, nil
		},
	}
	svc := cargroup.NewService(repo)

	got, err := svc.List(context.Background(), cargroup.ListFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(expected) {
		t.Errorf("expected %d items, got %d", len(expected), len(got))
	}
}

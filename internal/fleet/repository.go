package fleet

import (
	"context"

	"github.com/google/uuid"
)

// CarGroupRepository defines persistence operations for car groups.
type CarGroupRepository interface {
	// Create inserts a new car group and returns the persisted record.
	Create(ctx context.Context, input CreateCarGroupInput) (*CarGroup, error)

	// List returns car groups matching the provided filter.
	List(ctx context.Context, filter ListCarGroupsFilter) ([]*CarGroup, error)

	// GetByID returns the car group with the given ID.
	// It returns ErrNotFound when no matching record exists.
	GetByID(ctx context.Context, id uuid.UUID) (*CarGroup, error)

	// Update applies the provided changes to the car group identified by id.
	// It returns ErrNotFound when no matching record exists.
	Update(ctx context.Context, id uuid.UUID, input UpdateCarGroupInput) (*CarGroup, error)

	// Delete soft-deletes the car group identified by id.
	// It returns ErrNotFound when no matching record exists.
	// It returns ErrConflict when active vehicles reference this group.
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error

	// HasActiveVehicles reports whether any non-deleted vehicle references the car group.
	HasActiveVehicles(ctx context.Context, carGroupID uuid.UUID) (bool, error)
}

// VehicleRepository defines persistence operations for vehicles.
type VehicleRepository interface {
	// Create inserts a new vehicle and returns the persisted record.
	Create(ctx context.Context, input CreateVehicleInput) (*Vehicle, error)

	// List returns vehicles matching the provided filter (paginated).
	List(ctx context.Context, filter ListVehiclesFilter) ([]*Vehicle, int, error)

	// GetByID returns the vehicle with the given ID.
	// It returns ErrNotFound when no matching record exists.
	GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error)

	// Update applies the provided changes to the vehicle identified by id.
	// It returns ErrNotFound when no matching record exists.
	Update(ctx context.Context, id uuid.UUID, input UpdateVehicleInput) (*Vehicle, error)

	// UpdateDesignation changes the designation of the vehicle identified by id.
	// It returns ErrNotFound when no matching record exists.
	UpdateDesignation(ctx context.Context, id uuid.UUID, input UpdateDesignationInput) (*Vehicle, error)

	// Delete soft-deletes the vehicle identified by id.
	// It returns ErrNotFound when no matching record exists.
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

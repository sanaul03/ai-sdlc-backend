package fleet

import (
	"context"

	"github.com/google/uuid"
)

// CarGroupRepository defines the data access operations for car groups.
type CarGroupRepository interface {
	Create(ctx context.Context, group *CarGroup) error
	GetByID(ctx context.Context, id uuid.UUID) (*CarGroup, error)
	List(ctx context.Context) ([]*CarGroup, error)
	Update(ctx context.Context, group *CarGroup) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

// VehicleRepository defines the data access operations for vehicles.
type VehicleRepository interface {
	Create(ctx context.Context, vehicle *Vehicle) error
	GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error)
	List(ctx context.Context) ([]*Vehicle, error)
	Update(ctx context.Context, vehicle *Vehicle) error
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

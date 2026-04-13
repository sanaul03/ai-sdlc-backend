package repository

import (
	"context"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// Company defines the persistence contract for Company entities.
type Company interface {
	Create(ctx context.Context, req model.CreateCompanyRequest) (*model.Company, error)
	GetByID(ctx context.Context, id int64) (*model.Company, error)
	List(ctx context.Context) ([]model.Company, error)
	Update(ctx context.Context, id int64, req model.UpdateCompanyRequest) (*model.Company, error)
	Delete(ctx context.Context, id int64) error
}

// Depot defines the persistence contract for Depot entities.
type Depot interface {
	Create(ctx context.Context, req model.CreateDepotRequest) (*model.Depot, error)
	GetByID(ctx context.Context, id int64) (*model.Depot, error)
	ListByCompany(ctx context.Context, companyID int64) ([]model.Depot, error)
	Update(ctx context.Context, id int64, req model.UpdateDepotRequest) (*model.Depot, error)
	Delete(ctx context.Context, id int64) error
}

// VehicleCategory defines the persistence contract for VehicleCategory entities.
type VehicleCategory interface {
	Create(ctx context.Context, req model.CreateVehicleCategoryRequest) (*model.VehicleCategory, error)
	GetByID(ctx context.Context, id int64) (*model.VehicleCategory, error)
	List(ctx context.Context) ([]model.VehicleCategory, error)
	Update(ctx context.Context, id int64, req model.UpdateVehicleCategoryRequest) (*model.VehicleCategory, error)
	Delete(ctx context.Context, id int64) error
}

// VehicleType defines the persistence contract for VehicleType entities.
type VehicleType interface {
	Create(ctx context.Context, req model.CreateVehicleTypeRequest) (*model.VehicleType, error)
	GetByID(ctx context.Context, id int64) (*model.VehicleType, error)
	ListByCategory(ctx context.Context, categoryID int64) ([]model.VehicleType, error)
	Update(ctx context.Context, id int64, req model.UpdateVehicleTypeRequest) (*model.VehicleType, error)
	Delete(ctx context.Context, id int64) error
}

// Vehicle defines the persistence contract for Vehicle entities.
type Vehicle interface {
	Create(ctx context.Context, req model.CreateVehicleRequest) (*model.Vehicle, error)
	GetByID(ctx context.Context, id int64) (*model.Vehicle, error)
	ListByCompany(ctx context.Context, companyID int64) ([]model.Vehicle, error)
	Update(ctx context.Context, id int64, req model.UpdateVehicleRequest) (*model.Vehicle, error)
	Delete(ctx context.Context, id int64) error
}

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// VehicleRepository provides persistence operations for Vehicle entities.
type VehicleRepository struct {
	db *pgxpool.Pool
}

// NewVehicleRepository creates a new VehicleRepository.
func NewVehicleRepository(db *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// Create inserts a new vehicle and returns the created record.
func (r *VehicleRepository) Create(ctx context.Context, req model.CreateVehicleRequest) (*model.Vehicle, error) {
	const q = `
		INSERT INTO vehicles
			(company_id, depot_id, vehicle_type_id, registration_number,
			 chassis_number, engine_number, manufacture_year, color)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, company_id, depot_id, vehicle_type_id, registration_number,
			chassis_number, engine_number, manufacture_year, color, status,
			created_at, updated_at, deleted_at`

	v := &model.Vehicle{}
	err := r.db.QueryRow(ctx, q,
		req.CompanyID, req.DepotID, req.VehicleTypeID, req.RegistrationNumber,
		req.ChassisNumber, req.EngineNumber, req.ManufactureYear, req.Color,
	).Scan(&v.ID, &v.CompanyID, &v.DepotID, &v.VehicleTypeID, &v.RegistrationNumber,
		&v.ChassisNumber, &v.EngineNumber, &v.ManufactureYear, &v.Color,
		&v.Status, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// GetByID returns a single non-deleted vehicle by its primary key.
func (r *VehicleRepository) GetByID(ctx context.Context, id int64) (*model.Vehicle, error) {
	const q = `
		SELECT id, company_id, depot_id, vehicle_type_id, registration_number,
			chassis_number, engine_number, manufacture_year, color, status,
			created_at, updated_at, deleted_at
		FROM vehicles
		WHERE id = $1 AND deleted_at IS NULL`

	v := &model.Vehicle{}
	err := r.db.QueryRow(ctx, q, id).
		Scan(&v.ID, &v.CompanyID, &v.DepotID, &v.VehicleTypeID, &v.RegistrationNumber,
			&v.ChassisNumber, &v.EngineNumber, &v.ManufactureYear, &v.Color,
			&v.Status, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// ListByCompany returns all non-deleted vehicles belonging to the given company.
func (r *VehicleRepository) ListByCompany(ctx context.Context, companyID int64) ([]model.Vehicle, error) {
	const q = `
		SELECT id, company_id, depot_id, vehicle_type_id, registration_number,
			chassis_number, engine_number, manufacture_year, color, status,
			created_at, updated_at, deleted_at
		FROM vehicles
		WHERE company_id = $1 AND deleted_at IS NULL
		ORDER BY id`

	rows, err := r.db.Query(ctx, q, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Vehicle
	for rows.Next() {
		var v model.Vehicle
		if err := rows.Scan(&v.ID, &v.CompanyID, &v.DepotID, &v.VehicleTypeID, &v.RegistrationNumber,
			&v.ChassisNumber, &v.EngineNumber, &v.ManufactureYear, &v.Color,
			&v.Status, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

// Update modifies mutable fields of a vehicle and returns the updated record.
func (r *VehicleRepository) Update(ctx context.Context, id int64, req model.UpdateVehicleRequest) (*model.Vehicle, error) {
	const q = `
		UPDATE vehicles
		SET
			depot_id         = COALESCE($1, depot_id),
			vehicle_type_id  = COALESCE($2, vehicle_type_id),
			chassis_number   = COALESCE($3, chassis_number),
			engine_number    = COALESCE($4, engine_number),
			manufacture_year = COALESCE($5, manufacture_year),
			color            = COALESCE($6, color),
			status           = COALESCE($7, status),
			updated_at       = NOW()
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING id, company_id, depot_id, vehicle_type_id, registration_number,
			chassis_number, engine_number, manufacture_year, color, status,
			created_at, updated_at, deleted_at`

	v := &model.Vehicle{}
	err := r.db.QueryRow(ctx, q,
		req.DepotID, req.VehicleTypeID, req.ChassisNumber, req.EngineNumber,
		req.ManufactureYear, req.Color, req.Status, id,
	).Scan(&v.ID, &v.CompanyID, &v.DepotID, &v.VehicleTypeID, &v.RegistrationNumber,
		&v.ChassisNumber, &v.EngineNumber, &v.ManufactureYear, &v.Color,
		&v.Status, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Delete performs a soft-delete on a vehicle.
func (r *VehicleRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE vehicles SET deleted_at = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, time.Now(), id)
	return err
}

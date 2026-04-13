package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// VehicleTypeRepository provides persistence operations for VehicleType entities.
type VehicleTypeRepository struct {
	db *pgxpool.Pool
}

// NewVehicleTypeRepository creates a new VehicleTypeRepository.
func NewVehicleTypeRepository(db *pgxpool.Pool) *VehicleTypeRepository {
	return &VehicleTypeRepository{db: db}
}

// Create inserts a new vehicle type and returns the created record.
func (r *VehicleTypeRepository) Create(ctx context.Context, req model.CreateVehicleTypeRequest) (*model.VehicleType, error) {
	const q = `
		INSERT INTO vehicle_types (category_id, name, code, capacity, description)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, category_id, name, code, capacity, description, created_at, updated_at`

	vt := &model.VehicleType{}
	err := r.db.QueryRow(ctx, q,
		req.CategoryID, req.Name, req.Code, req.Capacity, req.Description,
	).Scan(&vt.ID, &vt.CategoryID, &vt.Name, &vt.Code, &vt.Capacity, &vt.Description, &vt.CreatedAt, &vt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vt, nil
}

// GetByID returns a single vehicle type by its primary key.
func (r *VehicleTypeRepository) GetByID(ctx context.Context, id int64) (*model.VehicleType, error) {
	const q = `
		SELECT id, category_id, name, code, capacity, description, created_at, updated_at
		FROM vehicle_types
		WHERE id = $1`

	vt := &model.VehicleType{}
	err := r.db.QueryRow(ctx, q, id).
		Scan(&vt.ID, &vt.CategoryID, &vt.Name, &vt.Code, &vt.Capacity, &vt.Description, &vt.CreatedAt, &vt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vt, nil
}

// ListByCategory returns all vehicle types belonging to the given category.
func (r *VehicleTypeRepository) ListByCategory(ctx context.Context, categoryID int64) ([]model.VehicleType, error) {
	const q = `
		SELECT id, category_id, name, code, capacity, description, created_at, updated_at
		FROM vehicle_types
		WHERE category_id = $1
		ORDER BY id`

	rows, err := r.db.Query(ctx, q, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.VehicleType
	for rows.Next() {
		var vt model.VehicleType
		if err := rows.Scan(&vt.ID, &vt.CategoryID, &vt.Name, &vt.Code, &vt.Capacity, &vt.Description, &vt.CreatedAt, &vt.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, vt)
	}
	return list, rows.Err()
}

// Update modifies mutable fields of a vehicle type and returns the updated record.
func (r *VehicleTypeRepository) Update(ctx context.Context, id int64, req model.UpdateVehicleTypeRequest) (*model.VehicleType, error) {
	const q = `
		UPDATE vehicle_types
		SET
			name        = COALESCE($1, name),
			capacity    = COALESCE($2, capacity),
			description = COALESCE($3, description),
			updated_at  = NOW()
		WHERE id = $4
		RETURNING id, category_id, name, code, capacity, description, created_at, updated_at`

	vt := &model.VehicleType{}
	err := r.db.QueryRow(ctx, q, req.Name, req.Capacity, req.Description, id).
		Scan(&vt.ID, &vt.CategoryID, &vt.Name, &vt.Code, &vt.Capacity, &vt.Description, &vt.CreatedAt, &vt.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vt, nil
}

// Delete removes a vehicle type permanently.
func (r *VehicleTypeRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM vehicle_types WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

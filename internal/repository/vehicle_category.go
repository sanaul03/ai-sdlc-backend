package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// VehicleCategoryRepository provides persistence operations for VehicleCategory entities.
type VehicleCategoryRepository struct {
	db *pgxpool.Pool
}

// NewVehicleCategoryRepository creates a new VehicleCategoryRepository.
func NewVehicleCategoryRepository(db *pgxpool.Pool) *VehicleCategoryRepository {
	return &VehicleCategoryRepository{db: db}
}

// Create inserts a new vehicle category and returns the created record.
func (r *VehicleCategoryRepository) Create(ctx context.Context, req model.CreateVehicleCategoryRequest) (*model.VehicleCategory, error) {
	const q = `
		INSERT INTO vehicle_categories (name, code, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, code, description, created_at, updated_at`

	vc := &model.VehicleCategory{}
	err := r.db.QueryRow(ctx, q, req.Name, req.Code, req.Description).
		Scan(&vc.ID, &vc.Name, &vc.Code, &vc.Description, &vc.CreatedAt, &vc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

// GetByID returns a single vehicle category by its primary key.
func (r *VehicleCategoryRepository) GetByID(ctx context.Context, id int64) (*model.VehicleCategory, error) {
	const q = `
		SELECT id, name, code, description, created_at, updated_at
		FROM vehicle_categories
		WHERE id = $1`

	vc := &model.VehicleCategory{}
	err := r.db.QueryRow(ctx, q, id).
		Scan(&vc.ID, &vc.Name, &vc.Code, &vc.Description, &vc.CreatedAt, &vc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

// List returns all vehicle categories.
func (r *VehicleCategoryRepository) List(ctx context.Context) ([]model.VehicleCategory, error) {
	const q = `
		SELECT id, name, code, description, created_at, updated_at
		FROM vehicle_categories
		ORDER BY id`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.VehicleCategory
	for rows.Next() {
		var vc model.VehicleCategory
		if err := rows.Scan(&vc.ID, &vc.Name, &vc.Code, &vc.Description, &vc.CreatedAt, &vc.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, vc)
	}
	return list, rows.Err()
}

// Update modifies mutable fields of a vehicle category and returns the updated record.
func (r *VehicleCategoryRepository) Update(ctx context.Context, id int64, req model.UpdateVehicleCategoryRequest) (*model.VehicleCategory, error) {
	const q = `
		UPDATE vehicle_categories
		SET
			name        = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at  = NOW()
		WHERE id = $3
		RETURNING id, name, code, description, created_at, updated_at`

	vc := &model.VehicleCategory{}
	err := r.db.QueryRow(ctx, q, req.Name, req.Description, id).
		Scan(&vc.ID, &vc.Name, &vc.Code, &vc.Description, &vc.CreatedAt, &vc.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return vc, nil
}

// Delete removes a vehicle category permanently.
func (r *VehicleCategoryRepository) Delete(ctx context.Context, id int64) error {
	const q = `DELETE FROM vehicle_categories WHERE id = $1`
	_, err := r.db.Exec(ctx, q, id)
	return err
}

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// DepotRepository provides persistence operations for Depot entities.
type DepotRepository struct {
	db *pgxpool.Pool
}

// NewDepotRepository creates a new DepotRepository.
func NewDepotRepository(db *pgxpool.Pool) *DepotRepository {
	return &DepotRepository{db: db}
}

// Create inserts a new depot and returns the created record.
func (r *DepotRepository) Create(ctx context.Context, req model.CreateDepotRequest) (*model.Depot, error) {
	const q = `
		INSERT INTO depots (company_id, name, code, address, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, company_id, name, code, address, latitude, longitude, status, created_at, updated_at, deleted_at`

	d := &model.Depot{}
	err := r.db.QueryRow(ctx, q,
		req.CompanyID, req.Name, req.Code, req.Address, req.Latitude, req.Longitude,
	).Scan(&d.ID, &d.CompanyID, &d.Name, &d.Code, &d.Address,
		&d.Latitude, &d.Longitude, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// GetByID returns a single non-deleted depot by its primary key.
func (r *DepotRepository) GetByID(ctx context.Context, id int64) (*model.Depot, error) {
	const q = `
		SELECT id, company_id, name, code, address, latitude, longitude, status, created_at, updated_at, deleted_at
		FROM depots
		WHERE id = $1 AND deleted_at IS NULL`

	d := &model.Depot{}
	err := r.db.QueryRow(ctx, q, id).
		Scan(&d.ID, &d.CompanyID, &d.Name, &d.Code, &d.Address,
			&d.Latitude, &d.Longitude, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// ListByCompany returns all non-deleted depots belonging to the given company.
func (r *DepotRepository) ListByCompany(ctx context.Context, companyID int64) ([]model.Depot, error) {
	const q = `
		SELECT id, company_id, name, code, address, latitude, longitude, status, created_at, updated_at, deleted_at
		FROM depots
		WHERE company_id = $1 AND deleted_at IS NULL
		ORDER BY id`

	rows, err := r.db.Query(ctx, q, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Depot
	for rows.Next() {
		var d model.Depot
		if err := rows.Scan(&d.ID, &d.CompanyID, &d.Name, &d.Code, &d.Address,
			&d.Latitude, &d.Longitude, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

// Update modifies mutable fields of a depot and returns the updated record.
func (r *DepotRepository) Update(ctx context.Context, id int64, req model.UpdateDepotRequest) (*model.Depot, error) {
	const q = `
		UPDATE depots
		SET
			name       = COALESCE($1, name),
			address    = COALESCE($2, address),
			latitude   = COALESCE($3, latitude),
			longitude  = COALESCE($4, longitude),
			status     = COALESCE($5, status),
			updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, company_id, name, code, address, latitude, longitude, status, created_at, updated_at, deleted_at`

	d := &model.Depot{}
	err := r.db.QueryRow(ctx, q,
		req.Name, req.Address, req.Latitude, req.Longitude, req.Status, id,
	).Scan(&d.ID, &d.CompanyID, &d.Name, &d.Code, &d.Address,
		&d.Latitude, &d.Longitude, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Delete performs a soft-delete on a depot.
func (r *DepotRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE depots SET deleted_at = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, time.Now(), id)
	return err
}

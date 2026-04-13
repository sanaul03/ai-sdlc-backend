package cargroup

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the persistence operations for car groups.
type Repository interface {
	Create(ctx context.Context, cg CarGroup) (CarGroup, error)
	List(ctx context.Context, filter ListFilter) ([]CarGroup, error)
	GetByID(ctx context.Context, id uuid.UUID) (CarGroup, error)
	Update(ctx context.Context, cg CarGroup) (CarGroup, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error
	HasActiveVehicles(ctx context.Context, id uuid.UUID) (bool, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

// NewRepository returns a PostgreSQL-backed Repository.
func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, cg CarGroup) (CarGroup, error) {
	const q = `
		INSERT INTO car_groups (id, name, description, size_category, created_at, updated_at, created_by, updated_by, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted`

	row := r.db.QueryRow(ctx, q,
		cg.ID, cg.Name, cg.Description, cg.SizeCategory,
		cg.CreatedAt, cg.UpdatedAt, cg.CreatedBy, cg.UpdatedBy, cg.Deleted,
	)

	var out CarGroup
	if err := scanCarGroup(row, &out); err != nil {
		return CarGroup{}, fmt.Errorf("cargroup.Repository.Create: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) List(ctx context.Context, filter ListFilter) ([]CarGroup, error) {
	q := `
		SELECT id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		WHERE deleted = $1`

	args := []any{filter.Deleted}

	if filter.Q != "" {
		q += fmt.Sprintf(` AND name ILIKE $%d`, len(args)+1)
		args = append(args, "%"+filter.Q+"%")
	}

	q += ` ORDER BY name ASC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("cargroup.Repository.List: %w", err)
	}
	defer rows.Close()

	var groups []CarGroup
	for rows.Next() {
		var cg CarGroup
		if err := scanCarGroupRow(rows, &cg); err != nil {
			return nil, fmt.Errorf("cargroup.Repository.List scan: %w", err)
		}
		groups = append(groups, cg)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("cargroup.Repository.List rows: %w", err)
	}
	return groups, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (CarGroup, error) {
	const q = `
		SELECT id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		WHERE id = $1`

	row := r.db.QueryRow(ctx, q, id)

	var cg CarGroup
	if err := scanCarGroup(row, &cg); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CarGroup{}, ErrNotFound
		}
		return CarGroup{}, fmt.Errorf("cargroup.Repository.GetByID: %w", err)
	}
	return cg, nil
}

func (r *postgresRepository) Update(ctx context.Context, cg CarGroup) (CarGroup, error) {
	const q = `
		UPDATE car_groups
		SET name = $2, description = $3, size_category = $4, updated_at = $5, updated_by = $6
		WHERE id = $1
		RETURNING id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted`

	row := r.db.QueryRow(ctx, q,
		cg.ID, cg.Name, cg.Description, cg.SizeCategory, cg.UpdatedAt, cg.UpdatedBy,
	)

	var out CarGroup
	if err := scanCarGroup(row, &out); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return CarGroup{}, ErrNotFound
		}
		return CarGroup{}, fmt.Errorf("cargroup.Repository.Update: %w", err)
	}
	return out, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	const q = `
		UPDATE car_groups
		SET deleted = true, deleted_at = now(), updated_at = now(), updated_by = $2
		WHERE id = $1 AND deleted = false`

	tag, err := r.db.Exec(ctx, q, id, deletedBy)
	if err != nil {
		return fmt.Errorf("cargroup.Repository.SoftDelete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresRepository) HasActiveVehicles(ctx context.Context, id uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS(SELECT 1 FROM vehicles WHERE car_group_id = $1 AND deleted = false)`

	var exists bool
	if err := r.db.QueryRow(ctx, q, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("cargroup.Repository.HasActiveVehicles: %w", err)
	}
	return exists, nil
}

// scanCarGroup scans a pgx.Row into a CarGroup.
func scanCarGroup(row pgx.Row, cg *CarGroup) error {
	return row.Scan(
		&cg.ID, &cg.Name, &cg.Description, &cg.SizeCategory,
		&cg.CreatedAt, &cg.UpdatedAt, &cg.DeletedAt,
		&cg.CreatedBy, &cg.UpdatedBy, &cg.Deleted,
	)
}

// scanCarGroupRow scans a pgx.Rows into a CarGroup.
func scanCarGroupRow(rows pgx.Rows, cg *CarGroup) error {
	return rows.Scan(
		&cg.ID, &cg.Name, &cg.Description, &cg.SizeCategory,
		&cg.CreatedAt, &cg.UpdatedAt, &cg.DeletedAt,
		&cg.CreatedBy, &cg.UpdatedBy, &cg.Deleted,
	)
}

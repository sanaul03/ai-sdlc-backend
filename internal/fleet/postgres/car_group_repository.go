package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// CarGroupRepository is a PostgreSQL-backed implementation of fleet.CarGroupRepository.
type CarGroupRepository struct {
	db *pgxpool.Pool
}

// NewCarGroupRepository constructs a new CarGroupRepository.
func NewCarGroupRepository(db *pgxpool.Pool) *CarGroupRepository {
	return &CarGroupRepository{db: db}
}

// Create inserts a new car group and returns the persisted record.
func (r *CarGroupRepository) Create(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error) {
	const q = `
		INSERT INTO car_groups (name, description, size_category, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, name, description, size_category,
		          created_at, updated_at, deleted_at, created_by, updated_by, deleted`

	row := r.db.QueryRow(ctx, q, input.Name, input.Description, input.SizeCategory, input.CreatedBy)
	return scanCarGroup(row)
}

// List returns car groups matching the provided filter.
func (r *CarGroupRepository) List(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error) {
	args := []any{}
	conditions := []string{}
	idx := 1

	if !filter.IncludeDeleted {
		conditions = append(conditions, fmt.Sprintf("deleted = $%d", idx))
		args = append(args, false)
		idx++
	}
	if filter.Q != nil && strings.TrimSpace(*filter.Q) != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+strings.TrimSpace(*filter.Q)+"%")
		idx++
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	q := fmt.Sprintf(`
		SELECT id, name, description, size_category,
		       created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		%s
		ORDER BY name`, where)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("car_group repo: list: %w", err)
	}
	defer rows.Close()

	var groups []*fleet.CarGroup
	for rows.Next() {
		g, err := scanCarGroupRow(rows)
		if err != nil {
			return nil, fmt.Errorf("car_group repo: list scan: %w", err)
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

// GetByID returns the car group with the given ID.
func (r *CarGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error) {
	const q = `
		SELECT id, name, description, size_category,
		       created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		WHERE id = $1`

	row := r.db.QueryRow(ctx, q, id)
	g, err := scanCarGroup(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fleet.ErrNotFound
	}
	return g, err
}

// Update applies changes to an existing car group.
func (r *CarGroupRepository) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error) {
	sets := []string{"updated_at = NOW()", fmt.Sprintf("updated_by = $%d", 1)}
	args := []any{input.UpdatedBy}
	idx := 2

	if input.Name != nil {
		sets = append(sets, fmt.Sprintf("name = $%d", idx))
		args = append(args, *input.Name)
		idx++
	}
	if input.Description != nil {
		sets = append(sets, fmt.Sprintf("description = $%d", idx))
		args = append(args, input.Description)
		idx++
	}
	if input.SizeCategory != nil {
		sets = append(sets, fmt.Sprintf("size_category = $%d", idx))
		args = append(args, input.SizeCategory)
		idx++
	}

	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE car_groups SET %s
		WHERE id = $%d AND deleted = false
		RETURNING id, name, description, size_category,
		          created_at, updated_at, deleted_at, created_by, updated_by, deleted`,
		strings.Join(sets, ", "), idx)

	row := r.db.QueryRow(ctx, q, args...)
	g, err := scanCarGroup(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fleet.ErrNotFound
	}
	return g, err
}

// Delete soft-deletes a car group.
func (r *CarGroupRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	const q = `
		UPDATE car_groups
		SET deleted = true, deleted_at = NOW(), updated_at = NOW(), updated_by = $2
		WHERE id = $1 AND deleted = false`

	tag, err := r.db.Exec(ctx, q, id, deletedBy)
	if err != nil {
		return fmt.Errorf("car_group repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fleet.ErrNotFound
	}
	return nil
}

// HasActiveVehicles reports whether any non-deleted vehicle references the car group.
func (r *CarGroupRepository) HasActiveVehicles(ctx context.Context, carGroupID uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS (SELECT 1 FROM vehicles WHERE car_group_id = $1 AND deleted = false)`
	var exists bool
	if err := r.db.QueryRow(ctx, q, carGroupID).Scan(&exists); err != nil {
		return false, fmt.Errorf("car_group repo: has_active_vehicles: %w", err)
	}
	return exists, nil
}

// scanCarGroup scans a single pgx.Row into a CarGroup.
func scanCarGroup(row pgx.Row) (*fleet.CarGroup, error) {
	var g fleet.CarGroup
	err := row.Scan(
		&g.ID, &g.Name, &g.Description, &g.SizeCategory,
		&g.CreatedAt, &g.UpdatedAt, &g.DeletedAt,
		&g.CreatedBy, &g.UpdatedBy, &g.Deleted,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// scanCarGroupRow scans a pgx.Rows into a CarGroup.
func scanCarGroupRow(rows pgx.Rows) (*fleet.CarGroup, error) {
	var g fleet.CarGroup
	err := rows.Scan(
		&g.ID, &g.Name, &g.Description, &g.SizeCategory,
		&g.CreatedAt, &g.UpdatedAt, &g.DeletedAt,
		&g.CreatedBy, &g.UpdatedBy, &g.Deleted,
	)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

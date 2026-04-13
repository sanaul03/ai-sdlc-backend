package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// CompanyRepository provides persistence operations for Company entities.
type CompanyRepository struct {
	db *pgxpool.Pool
}

// NewCompanyRepository creates a new CompanyRepository.
func NewCompanyRepository(db *pgxpool.Pool) *CompanyRepository {
	return &CompanyRepository{db: db}
}

// Create inserts a new company and returns the created record.
func (r *CompanyRepository) Create(ctx context.Context, req model.CreateCompanyRequest) (*model.Company, error) {
	const q = `
		INSERT INTO companies (name, code, address, phone, email)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, code, address, phone, email, status, created_at, updated_at, deleted_at`

	c := &model.Company{}
	err := r.db.QueryRow(ctx, q,
		req.Name, req.Code, req.Address, req.Phone, req.Email,
	).Scan(&c.ID, &c.Name, &c.Code, &c.Address, &c.Phone, &c.Email,
		&c.Status, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// GetByID returns a single non-deleted company by its primary key.
func (r *CompanyRepository) GetByID(ctx context.Context, id int64) (*model.Company, error) {
	const q = `
		SELECT id, name, code, address, phone, email, status, created_at, updated_at, deleted_at
		FROM companies
		WHERE id = $1 AND deleted_at IS NULL`

	c := &model.Company{}
	err := r.db.QueryRow(ctx, q, id).
		Scan(&c.ID, &c.Name, &c.Code, &c.Address, &c.Phone, &c.Email,
			&c.Status, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// List returns all non-deleted companies.
func (r *CompanyRepository) List(ctx context.Context) ([]model.Company, error) {
	const q = `
		SELECT id, name, code, address, phone, email, status, created_at, updated_at, deleted_at
		FROM companies
		WHERE deleted_at IS NULL
		ORDER BY id`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Company
	for rows.Next() {
		var c model.Company
		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &c.Address, &c.Phone, &c.Email,
			&c.Status, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// Update modifies mutable fields of a company and returns the updated record.
func (r *CompanyRepository) Update(ctx context.Context, id int64, req model.UpdateCompanyRequest) (*model.Company, error) {
	const q = `
		UPDATE companies
		SET
			name       = COALESCE($1, name),
			address    = COALESCE($2, address),
			phone      = COALESCE($3, phone),
			email      = COALESCE($4, email),
			status     = COALESCE($5, status),
			updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, name, code, address, phone, email, status, created_at, updated_at, deleted_at`

	c := &model.Company{}
	err := r.db.QueryRow(ctx, q,
		req.Name, req.Address, req.Phone, req.Email, req.Status, id,
	).Scan(&c.ID, &c.Name, &c.Code, &c.Address, &c.Phone, &c.Email,
		&c.Status, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Delete performs a soft-delete on a company.
func (r *CompanyRepository) Delete(ctx context.Context, id int64) error {
	const q = `UPDATE companies SET deleted_at = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.Exec(ctx, q, time.Now(), id)
	return err
}

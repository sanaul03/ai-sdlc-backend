package vehicle

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the persistence operations for vehicles.
type Repository interface {
	Create(ctx context.Context, v Vehicle) (Vehicle, error)
	List(ctx context.Context, filter ListFilter) (Page, error)
	GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error)
	Update(ctx context.Context, v Vehicle) (Vehicle, error)
	UpdateDesignation(ctx context.Context, id uuid.UUID, designation, updatedBy string) (Vehicle, error)
	SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

// NewRepository returns a PostgreSQL-backed Repository.
func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

const vehicleColumns = `
	id, car_group_id, branch_id, vin, licence_plate, brand, model, year, colour,
	fuel_type, transmission_type, current_mileage, status, designation,
	acquisition_date, ownership_type, lease_details, insurance_policy_number,
	insurance_expiry_date, registration_expiry_date, last_inspection_date,
	next_inspection_due_date, notes, created_at, updated_at, deleted_at,
	created_by, updated_by, deleted`

func (r *postgresRepository) Create(ctx context.Context, v Vehicle) (Vehicle, error) {
	q := fmt.Sprintf(`
		INSERT INTO vehicles (
			id, car_group_id, branch_id, vin, licence_plate, brand, model, year, colour,
			fuel_type, transmission_type, current_mileage, status, designation,
			acquisition_date, ownership_type, lease_details, insurance_policy_number,
			insurance_expiry_date, registration_expiry_date, last_inspection_date,
			next_inspection_due_date, notes, created_at, updated_at, created_by, updated_by, deleted
		) VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28
		) RETURNING %s`, vehicleColumns)

	row := r.db.QueryRow(ctx, q,
		v.ID, v.CarGroupID, v.BranchID, v.VIN, v.LicencePlate,
		v.Brand, v.Model, v.Year, v.Colour,
		v.FuelType, v.TransmissionType, v.CurrentMileage, v.Status, v.Designation,
		v.AcquisitionDate, v.OwnershipType, v.LeaseDetails, v.InsurancePolicyNumber,
		v.InsuranceExpiryDate, v.RegistrationExpiryDate, v.LastInspectionDate,
		v.NextInspectionDueDate, v.Notes, v.CreatedAt, v.UpdatedAt, v.CreatedBy, v.UpdatedBy, v.Deleted,
	)

	var out Vehicle
	if err := scanVehicle(row, &out); err != nil {
		return Vehicle{}, fmt.Errorf("vehicle.Repository.Create: %w", err)
	}
	setExpiryWarning(&out)
	return out, nil
}

func (r *postgresRepository) List(ctx context.Context, filter ListFilter) (Page, error) {
	where := []string{"v.deleted = false"}
	args := []any{}

	if filter.CarGroupID != nil {
		args = append(args, *filter.CarGroupID)
		where = append(where, fmt.Sprintf("v.car_group_id = $%d", len(args)))
	}
	if filter.BranchID != nil {
		args = append(args, *filter.BranchID)
		where = append(where, fmt.Sprintf("v.branch_id = $%d", len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		where = append(where, fmt.Sprintf("v.status = $%d", len(args)))
	}
	if filter.Designation != "" {
		args = append(args, filter.Designation)
		where = append(where, fmt.Sprintf("v.designation = $%d", len(args)))
	}
	if filter.FuelType != "" {
		args = append(args, filter.FuelType)
		where = append(where, fmt.Sprintf("v.fuel_type = $%d", len(args)))
	}
	if filter.TransmissionType != "" {
		args = append(args, filter.TransmissionType)
		where = append(where, fmt.Sprintf("v.transmission_type = $%d", len(args)))
	}
	if filter.ExpiryWarning {
		where = append(where, `(
			(v.insurance_expiry_date IS NOT NULL AND v.insurance_expiry_date <= now() + INTERVAL '30 days') OR
			(v.registration_expiry_date IS NOT NULL AND v.registration_expiry_date <= now() + INTERVAL '30 days')
		)`)
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQ := fmt.Sprintf(`SELECT COUNT(*) FROM vehicles v WHERE %s`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return Page{}, fmt.Errorf("vehicle.Repository.List count: %w", err)
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	args = append(args, pageSize, offset)
	dataQ := fmt.Sprintf(`
		SELECT %s FROM vehicles v
		WHERE %s
		ORDER BY v.created_at DESC
		LIMIT $%d OFFSET $%d`,
		vehicleColumns, whereClause, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, dataQ, args...)
	if err != nil {
		return Page{}, fmt.Errorf("vehicle.Repository.List: %w", err)
	}
	defer rows.Close()

	items := make([]Vehicle, 0)
	for rows.Next() {
		var v Vehicle
		if err := scanVehicleRow(rows, &v); err != nil {
			return Page{}, fmt.Errorf("vehicle.Repository.List scan: %w", err)
		}
		setExpiryWarning(&v)
		items = append(items, v)
	}
	if err := rows.Err(); err != nil {
		return Page{}, fmt.Errorf("vehicle.Repository.List rows: %w", err)
	}

	return Page{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id uuid.UUID) (Vehicle, error) {
	q := fmt.Sprintf(`SELECT %s FROM vehicles v WHERE v.id = $1`, vehicleColumns)

	row := r.db.QueryRow(ctx, q, id)
	var v Vehicle
	if err := scanVehicle(row, &v); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vehicle{}, ErrNotFound
		}
		return Vehicle{}, fmt.Errorf("vehicle.Repository.GetByID: %w", err)
	}
	setExpiryWarning(&v)
	return v, nil
}

func (r *postgresRepository) Update(ctx context.Context, v Vehicle) (Vehicle, error) {
	q := fmt.Sprintf(`
		UPDATE vehicles SET
			car_group_id = $2, branch_id = $3, vin = $4, licence_plate = $5,
			brand = $6, model = $7, year = $8, colour = $9,
			fuel_type = $10, transmission_type = $11, current_mileage = $12,
			designation = $13, acquisition_date = $14, ownership_type = $15,
			lease_details = $16, insurance_policy_number = $17,
			insurance_expiry_date = $18, registration_expiry_date = $19,
			last_inspection_date = $20, next_inspection_due_date = $21,
			notes = $22, updated_at = $23, updated_by = $24
		WHERE id = $1
		RETURNING %s`, vehicleColumns)

	row := r.db.QueryRow(ctx, q,
		v.ID, v.CarGroupID, v.BranchID, v.VIN, v.LicencePlate,
		v.Brand, v.Model, v.Year, v.Colour,
		v.FuelType, v.TransmissionType, v.CurrentMileage,
		v.Designation, v.AcquisitionDate, v.OwnershipType,
		v.LeaseDetails, v.InsurancePolicyNumber,
		v.InsuranceExpiryDate, v.RegistrationExpiryDate,
		v.LastInspectionDate, v.NextInspectionDueDate,
		v.Notes, v.UpdatedAt, v.UpdatedBy,
	)

	var out Vehicle
	if err := scanVehicle(row, &out); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vehicle{}, ErrNotFound
		}
		return Vehicle{}, fmt.Errorf("vehicle.Repository.Update: %w", err)
	}
	setExpiryWarning(&out)
	return out, nil
}

func (r *postgresRepository) UpdateDesignation(ctx context.Context, id uuid.UUID, designation, updatedBy string) (Vehicle, error) {
	q := fmt.Sprintf(`
		UPDATE vehicles SET
			designation = $2, updated_at = now(), updated_by = $3
		WHERE id = $1
		RETURNING %s`, vehicleColumns)

	row := r.db.QueryRow(ctx, q, id, designation, updatedBy)
	var out Vehicle
	if err := scanVehicle(row, &out); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Vehicle{}, ErrNotFound
		}
		return Vehicle{}, fmt.Errorf("vehicle.Repository.UpdateDesignation: %w", err)
	}
	setExpiryWarning(&out)
	return out, nil
}

func (r *postgresRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	const q = `
		UPDATE vehicles
		SET deleted = true, deleted_at = now(), updated_at = now(), updated_by = $2
		WHERE id = $1 AND deleted = false`

	tag, err := r.db.Exec(ctx, q, id, deletedBy)
	if err != nil {
		return fmt.Errorf("vehicle.Repository.SoftDelete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// scanVehicle scans a pgx.Row into a Vehicle.
func scanVehicle(row pgx.Row, v *Vehicle) error {
	return row.Scan(
		&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
		&v.Brand, &v.Model, &v.Year, &v.Colour,
		&v.FuelType, &v.TransmissionType, &v.CurrentMileage, &v.Status, &v.Designation,
		&v.AcquisitionDate, &v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
		&v.InsuranceExpiryDate, &v.RegistrationExpiryDate, &v.LastInspectionDate,
		&v.NextInspectionDueDate, &v.Notes, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
		&v.CreatedBy, &v.UpdatedBy, &v.Deleted,
	)
}

// scanVehicleRow scans a pgx.Rows into a Vehicle.
func scanVehicleRow(rows pgx.Rows, v *Vehicle) error {
	return rows.Scan(
		&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
		&v.Brand, &v.Model, &v.Year, &v.Colour,
		&v.FuelType, &v.TransmissionType, &v.CurrentMileage, &v.Status, &v.Designation,
		&v.AcquisitionDate, &v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
		&v.InsuranceExpiryDate, &v.RegistrationExpiryDate, &v.LastInspectionDate,
		&v.NextInspectionDueDate, &v.Notes, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
		&v.CreatedBy, &v.UpdatedBy, &v.Deleted,
	)
}

// setExpiryWarning sets ExpiryWarning=true if insurance or registration
// expires within 30 days from now.
func setExpiryWarning(v *Vehicle) {
	threshold := time.Now().UTC().Add(30 * 24 * time.Hour)
	if v.InsuranceExpiryDate != nil && v.InsuranceExpiryDate.Before(threshold) {
		v.ExpiryWarning = true
		return
	}
	if v.RegistrationExpiryDate != nil && v.RegistrationExpiryDate.Before(threshold) {
		v.ExpiryWarning = true
	}
}

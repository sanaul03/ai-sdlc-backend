package fleet

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgCarGroupRepository is the PostgreSQL implementation of CarGroupRepository.
type pgCarGroupRepository struct {
	db *pgxpool.Pool
}

// NewPgCarGroupRepository creates a new PostgreSQL-backed CarGroupRepository.
func NewPgCarGroupRepository(db *pgxpool.Pool) CarGroupRepository {
	return &pgCarGroupRepository{db: db}
}

func (r *pgCarGroupRepository) Create(ctx context.Context, group *CarGroup) error {
	query := `
		INSERT INTO car_groups (id, name, description, size_category, created_at, updated_at, created_by, updated_by, deleted)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		group.ID, group.Name, group.Description, group.SizeCategory,
		group.CreatedAt, group.UpdatedAt, group.CreatedBy, group.UpdatedBy, group.Deleted,
	)
	if err != nil {
		return fmt.Errorf("create car group: %w", err)
	}
	return nil
}

func (r *pgCarGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*CarGroup, error) {
	query := `
		SELECT id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		WHERE id = $1 AND deleted = false
	`
	row := r.db.QueryRow(ctx, query, id)
	group := &CarGroup{}
	err := row.Scan(
		&group.ID, &group.Name, &group.Description, &group.SizeCategory,
		&group.CreatedAt, &group.UpdatedAt, &group.DeletedAt,
		&group.CreatedBy, &group.UpdatedBy, &group.Deleted,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get car group by id: %w", err)
	}
	return group, nil
}

func (r *pgCarGroupRepository) List(ctx context.Context) ([]*CarGroup, error) {
	query := `
		SELECT id, name, description, size_category, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM car_groups
		WHERE deleted = false
		ORDER BY name ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list car groups: %w", err)
	}
	defer rows.Close()

	var groups []*CarGroup
	for rows.Next() {
		group := &CarGroup{}
		if err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.SizeCategory,
			&group.CreatedAt, &group.UpdatedAt, &group.DeletedAt,
			&group.CreatedBy, &group.UpdatedBy, &group.Deleted,
		); err != nil {
			return nil, fmt.Errorf("scan car group: %w", err)
		}
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate car groups: %w", err)
	}
	return groups, nil
}

func (r *pgCarGroupRepository) Update(ctx context.Context, group *CarGroup) error {
	query := `
		UPDATE car_groups
		SET name = $1, description = $2, size_category = $3, updated_at = $4, updated_by = $5
		WHERE id = $6 AND deleted = false
	`
	result, err := r.db.Exec(ctx, query,
		group.Name, group.Description, group.SizeCategory,
		group.UpdatedAt, group.UpdatedBy, group.ID,
	)
	if err != nil {
		return fmt.Errorf("update car group: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgCarGroupRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	now := time.Now().UTC()
	query := `
		UPDATE car_groups
		SET deleted = true, deleted_at = $1, updated_at = $2, updated_by = $3
		WHERE id = $4 AND deleted = false
	`
	result, err := r.db.Exec(ctx, query, now, now, deletedBy, id)
	if err != nil {
		return fmt.Errorf("delete car group: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// pgVehicleRepository is the PostgreSQL implementation of VehicleRepository.
type pgVehicleRepository struct {
	db *pgxpool.Pool
}

// NewPgVehicleRepository creates a new PostgreSQL-backed VehicleRepository.
func NewPgVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &pgVehicleRepository{db: db}
}

func (r *pgVehicleRepository) Create(ctx context.Context, vehicle *Vehicle) error {
	query := `
		INSERT INTO vehicles (
			id, car_group_id, branch_id, vin, licence_plate, brand, model, year, colour,
			fuel_type, transmission_type, current_mileage, status, designation,
			acquisition_date, ownership_type, lease_details, insurance_policy_number,
			insurance_expiry_date, registration_expiry_date, last_inspection_date,
			next_inspection_due_date, notes, created_at, updated_at, created_by, updated_by, deleted
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14,
			$15, $16, $17, $18,
			$19, $20, $21,
			$22, $23, $24, $25, $26, $27, $28
		)
	`
	_, err := r.db.Exec(ctx, query,
		vehicle.ID, vehicle.CarGroupID, vehicle.BranchID, vehicle.VIN, vehicle.LicencePlate,
		vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Colour,
		vehicle.FuelType, vehicle.TransmissionType, vehicle.CurrentMileage, vehicle.Status, vehicle.Designation,
		vehicle.AcquisitionDate, vehicle.OwnershipType, vehicle.LeaseDetails, vehicle.InsurancePolicyNumber,
		vehicle.InsuranceExpiryDate, vehicle.RegistrationExpiryDate, vehicle.LastInspectionDate,
		vehicle.NextInspectionDueDate, vehicle.Notes, vehicle.CreatedAt, vehicle.UpdatedAt,
		vehicle.CreatedBy, vehicle.UpdatedBy, vehicle.Deleted,
	)
	if err != nil {
		return fmt.Errorf("create vehicle: %w", err)
	}
	return nil
}

func (r *pgVehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*Vehicle, error) {
	query := `
		SELECT
			id, car_group_id, branch_id, vin, licence_plate, brand, model, year, colour,
			fuel_type, transmission_type, current_mileage, status, designation,
			acquisition_date, ownership_type, lease_details, insurance_policy_number,
			insurance_expiry_date, registration_expiry_date, last_inspection_date,
			next_inspection_due_date, notes, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM vehicles
		WHERE id = $1 AND deleted = false
	`
	row := r.db.QueryRow(ctx, query, id)
	v := &Vehicle{}
	err := row.Scan(
		&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
		&v.Brand, &v.Model, &v.Year, &v.Colour,
		&v.FuelType, &v.TransmissionType, &v.CurrentMileage, &v.Status, &v.Designation,
		&v.AcquisitionDate, &v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
		&v.InsuranceExpiryDate, &v.RegistrationExpiryDate, &v.LastInspectionDate,
		&v.NextInspectionDueDate, &v.Notes, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
		&v.CreatedBy, &v.UpdatedBy, &v.Deleted,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get vehicle by id: %w", err)
	}
	return v, nil
}

func (r *pgVehicleRepository) List(ctx context.Context) ([]*Vehicle, error) {
	query := `
		SELECT
			id, car_group_id, branch_id, vin, licence_plate, brand, model, year, colour,
			fuel_type, transmission_type, current_mileage, status, designation,
			acquisition_date, ownership_type, lease_details, insurance_policy_number,
			insurance_expiry_date, registration_expiry_date, last_inspection_date,
			next_inspection_due_date, notes, created_at, updated_at, deleted_at, created_by, updated_by, deleted
		FROM vehicles
		WHERE deleted = false
		ORDER BY brand ASC, model ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []*Vehicle
	for rows.Next() {
		v := &Vehicle{}
		if err := rows.Scan(
			&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
			&v.Brand, &v.Model, &v.Year, &v.Colour,
			&v.FuelType, &v.TransmissionType, &v.CurrentMileage, &v.Status, &v.Designation,
			&v.AcquisitionDate, &v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
			&v.InsuranceExpiryDate, &v.RegistrationExpiryDate, &v.LastInspectionDate,
			&v.NextInspectionDueDate, &v.Notes, &v.CreatedAt, &v.UpdatedAt, &v.DeletedAt,
			&v.CreatedBy, &v.UpdatedBy, &v.Deleted,
		); err != nil {
			return nil, fmt.Errorf("scan vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate vehicles: %w", err)
	}
	return vehicles, nil
}

func (r *pgVehicleRepository) Update(ctx context.Context, vehicle *Vehicle) error {
	query := `
		UPDATE vehicles
		SET car_group_id = $1, branch_id = $2, vin = $3, licence_plate = $4, brand = $5,
		    model = $6, year = $7, colour = $8, fuel_type = $9, transmission_type = $10,
		    current_mileage = $11, status = $12, designation = $13, acquisition_date = $14,
		    ownership_type = $15, lease_details = $16, insurance_policy_number = $17,
		    insurance_expiry_date = $18, registration_expiry_date = $19,
		    last_inspection_date = $20, next_inspection_due_date = $21, notes = $22,
		    updated_at = $23, updated_by = $24
		WHERE id = $25 AND deleted = false
	`
	result, err := r.db.Exec(ctx, query,
		vehicle.CarGroupID, vehicle.BranchID, vehicle.VIN, vehicle.LicencePlate,
		vehicle.Brand, vehicle.Model, vehicle.Year, vehicle.Colour,
		vehicle.FuelType, vehicle.TransmissionType, vehicle.CurrentMileage,
		vehicle.Status, vehicle.Designation, vehicle.AcquisitionDate,
		vehicle.OwnershipType, vehicle.LeaseDetails, vehicle.InsurancePolicyNumber,
		vehicle.InsuranceExpiryDate, vehicle.RegistrationExpiryDate,
		vehicle.LastInspectionDate, vehicle.NextInspectionDueDate, vehicle.Notes,
		vehicle.UpdatedAt, vehicle.UpdatedBy, vehicle.ID,
	)
	if err != nil {
		return fmt.Errorf("update vehicle: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgVehicleRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	now := time.Now().UTC()
	query := `
		UPDATE vehicles
		SET deleted = true, deleted_at = $1, updated_at = $2, updated_by = $3
		WHERE id = $4 AND deleted = false
	`
	result, err := r.db.Exec(ctx, query, now, now, deletedBy, id)
	if err != nil {
		return fmt.Errorf("delete vehicle: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

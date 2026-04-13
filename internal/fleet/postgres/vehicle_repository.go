package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

const expiryWarningDays = 30

// VehicleRepository is a PostgreSQL-backed implementation of fleet.VehicleRepository.
type VehicleRepository struct {
	db *pgxpool.Pool
}

// NewVehicleRepository constructs a new VehicleRepository.
func NewVehicleRepository(db *pgxpool.Pool) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// Create inserts a new vehicle and returns the persisted record.
func (r *VehicleRepository) Create(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
	const q = `
		INSERT INTO vehicles (
			car_group_id, branch_id, vin, licence_plate, brand, model, year,
			colour, fuel_type, transmission_type, current_mileage,
			status, designation, acquisition_date, ownership_type, lease_details,
			insurance_policy_number, insurance_expiry_date, registration_expiry_date,
			last_inspection_date, next_inspection_due_date, notes,
			created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11,
			$12, $13, $14, $15, $16,
			$17, $18, $19,
			$20, $21, $22,
			$23, $23
		)
		RETURNING ` + vehicleColumns

	row := r.db.QueryRow(ctx, q,
		input.CarGroupID, input.BranchID, input.VIN, input.LicencePlate,
		input.Brand, input.Model, input.Year,
		input.Colour, string(input.FuelType), string(input.TransmissionType), input.CurrentMileage,
		string(input.Status), string(input.Designation), input.AcquisitionDate, input.OwnershipType, input.LeaseDetails,
		input.InsurancePolicyNumber, input.InsuranceExpiryDate, input.RegistrationExpiryDate,
		input.LastInspectionDate, input.NextInspectionDueDate, input.Notes,
		input.CreatedBy,
	)
	return scanVehicle(row)
}

// List returns vehicles matching the filter (paginated).
func (r *VehicleRepository) List(ctx context.Context, filter fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error) {
	args := []any{}
	conditions := []string{"v.deleted = false"}
	idx := 1

	if filter.CarGroupID != nil {
		conditions = append(conditions, fmt.Sprintf("v.car_group_id = $%d", idx))
		args = append(args, *filter.CarGroupID)
		idx++
	}
	if filter.BranchID != nil {
		conditions = append(conditions, fmt.Sprintf("v.branch_id = $%d", idx))
		args = append(args, *filter.BranchID)
		idx++
	}
	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("v.status = $%d", idx))
		args = append(args, string(*filter.Status))
		idx++
	}
	if filter.Designation != nil {
		conditions = append(conditions, fmt.Sprintf("v.designation = $%d", idx))
		args = append(args, string(*filter.Designation))
		idx++
	}
	if filter.FuelType != nil {
		conditions = append(conditions, fmt.Sprintf("v.fuel_type = $%d", idx))
		args = append(args, string(*filter.FuelType))
		idx++
	}
	if filter.TransmissionType != nil {
		conditions = append(conditions, fmt.Sprintf("v.transmission_type = $%d", idx))
		args = append(args, string(*filter.TransmissionType))
		idx++
	}
	if filter.ExpiryWarning != nil && *filter.ExpiryWarning {
		threshold := time.Now().AddDate(0, 0, expiryWarningDays)
		conditions = append(conditions, fmt.Sprintf(
			"(v.insurance_expiry_date <= $%d OR v.registration_expiry_date <= $%d)", idx, idx))
		args = append(args, threshold)
		idx++
	}

	where := "WHERE " + strings.Join(conditions, " AND ")

	// Count query
	countQ := fmt.Sprintf("SELECT COUNT(*) FROM vehicles v %s", where)
	var total int
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("vehicle repo: list count: %w", err)
	}

	// Data query with pagination
	page := filter.Page
	pageSize := filter.PageSize
	offset := (page - 1) * pageSize
	args = append(args, pageSize, offset)
	dataQ := fmt.Sprintf(`
		SELECT %s FROM vehicles v %s
		ORDER BY v.created_at DESC
		LIMIT $%d OFFSET $%d`, vehicleColumns, where, idx, idx+1)

	rows, err := r.db.Query(ctx, dataQ, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("vehicle repo: list query: %w", err)
	}
	defer rows.Close()

	var vehicles []*fleet.Vehicle
	for rows.Next() {
		v, err := scanVehicleRow(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("vehicle repo: list scan: %w", err)
		}
		setExpiryWarning(v)
		vehicles = append(vehicles, v)
	}
	return vehicles, total, rows.Err()
}

// GetByID returns the vehicle with the given ID.
func (r *VehicleRepository) GetByID(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error) {
	q := fmt.Sprintf(`SELECT %s FROM vehicles v WHERE v.id = $1`, vehicleColumns)
	row := r.db.QueryRow(ctx, q, id)
	v, err := scanVehicle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fleet.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	setExpiryWarning(v)
	return v, nil
}

// Update applies changes to an existing vehicle.
func (r *VehicleRepository) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error) {
	sets := []string{"updated_at = NOW()", fmt.Sprintf("updated_by = $%d", 1)}
	args := []any{input.UpdatedBy}
	idx := 2

	if input.CarGroupID != nil {
		sets = append(sets, fmt.Sprintf("car_group_id = $%d", idx))
		args = append(args, *input.CarGroupID)
		idx++
	}
	if input.BranchID != nil {
		sets = append(sets, fmt.Sprintf("branch_id = $%d", idx))
		args = append(args, *input.BranchID)
		idx++
	}
	if input.VIN != nil {
		sets = append(sets, fmt.Sprintf("vin = $%d", idx))
		args = append(args, *input.VIN)
		idx++
	}
	if input.LicencePlate != nil {
		sets = append(sets, fmt.Sprintf("licence_plate = $%d", idx))
		args = append(args, *input.LicencePlate)
		idx++
	}
	if input.Brand != nil {
		sets = append(sets, fmt.Sprintf("brand = $%d", idx))
		args = append(args, *input.Brand)
		idx++
	}
	if input.Model != nil {
		sets = append(sets, fmt.Sprintf("model = $%d", idx))
		args = append(args, *input.Model)
		idx++
	}
	if input.Year != nil {
		sets = append(sets, fmt.Sprintf("year = $%d", idx))
		args = append(args, *input.Year)
		idx++
	}
	if input.Colour != nil {
		sets = append(sets, fmt.Sprintf("colour = $%d", idx))
		args = append(args, input.Colour)
		idx++
	}
	if input.FuelType != nil {
		sets = append(sets, fmt.Sprintf("fuel_type = $%d", idx))
		args = append(args, string(*input.FuelType))
		idx++
	}
	if input.TransmissionType != nil {
		sets = append(sets, fmt.Sprintf("transmission_type = $%d", idx))
		args = append(args, string(*input.TransmissionType))
		idx++
	}
	if input.CurrentMileage != nil {
		sets = append(sets, fmt.Sprintf("current_mileage = $%d", idx))
		args = append(args, *input.CurrentMileage)
		idx++
	}
	if input.Designation != nil {
		sets = append(sets, fmt.Sprintf("designation = $%d", idx))
		args = append(args, string(*input.Designation))
		idx++
	}
	if input.AcquisitionDate != nil {
		sets = append(sets, fmt.Sprintf("acquisition_date = $%d", idx))
		args = append(args, *input.AcquisitionDate)
		idx++
	}
	if input.OwnershipType != nil {
		sets = append(sets, fmt.Sprintf("ownership_type = $%d", idx))
		args = append(args, *input.OwnershipType)
		idx++
	}
	if input.LeaseDetails != nil {
		sets = append(sets, fmt.Sprintf("lease_details = $%d", idx))
		args = append(args, input.LeaseDetails)
		idx++
	}
	if input.InsurancePolicyNumber != nil {
		sets = append(sets, fmt.Sprintf("insurance_policy_number = $%d", idx))
		args = append(args, input.InsurancePolicyNumber)
		idx++
	}
	if input.InsuranceExpiryDate != nil {
		sets = append(sets, fmt.Sprintf("insurance_expiry_date = $%d", idx))
		args = append(args, input.InsuranceExpiryDate)
		idx++
	}
	if input.RegistrationExpiryDate != nil {
		sets = append(sets, fmt.Sprintf("registration_expiry_date = $%d", idx))
		args = append(args, input.RegistrationExpiryDate)
		idx++
	}
	if input.LastInspectionDate != nil {
		sets = append(sets, fmt.Sprintf("last_inspection_date = $%d", idx))
		args = append(args, input.LastInspectionDate)
		idx++
	}
	if input.NextInspectionDueDate != nil {
		sets = append(sets, fmt.Sprintf("next_inspection_due_date = $%d", idx))
		args = append(args, input.NextInspectionDueDate)
		idx++
	}
	if input.Notes != nil {
		sets = append(sets, fmt.Sprintf("notes = $%d", idx))
		args = append(args, input.Notes)
		idx++
	}

	args = append(args, id)
	q := fmt.Sprintf(`
		UPDATE vehicles SET %s
		WHERE id = $%d AND deleted = false
		RETURNING %s`, strings.Join(sets, ", "), idx, vehicleColumns)

	row := r.db.QueryRow(ctx, q, args...)
	v, err := scanVehicle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fleet.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	setExpiryWarning(v)
	return v, nil
}

// UpdateDesignation changes the designation of a vehicle.
func (r *VehicleRepository) UpdateDesignation(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error) {
	q := fmt.Sprintf(`
		UPDATE vehicles
		SET designation = $1, updated_at = NOW(), updated_by = $2
		WHERE id = $3 AND deleted = false
		RETURNING %s`, vehicleColumns)

	row := r.db.QueryRow(ctx, q, string(input.Designation), input.UpdatedBy, id)
	v, err := scanVehicle(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fleet.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	setExpiryWarning(v)
	return v, nil
}

// Delete soft-deletes a vehicle.
func (r *VehicleRepository) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	const q = `
		UPDATE vehicles
		SET deleted = true, deleted_at = NOW(), updated_at = NOW(), updated_by = $2
		WHERE id = $1 AND deleted = false`

	tag, err := r.db.Exec(ctx, q, id, deletedBy)
	if err != nil {
		return fmt.Errorf("vehicle repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fleet.ErrNotFound
	}
	return nil
}

// vehicleColumns is the shared SELECT column list for the vehicles table (aliased as v).
const vehicleColumns = `
	v.id, v.car_group_id, v.branch_id, v.vin, v.licence_plate,
	v.brand, v.model, v.year, v.colour, v.fuel_type, v.transmission_type,
	v.current_mileage, v.status, v.designation, v.acquisition_date,
	v.ownership_type, v.lease_details, v.insurance_policy_number,
	v.insurance_expiry_date, v.registration_expiry_date,
	v.last_inspection_date, v.next_inspection_due_date, v.notes,
	v.created_at, v.updated_at, v.deleted_at, v.created_by, v.updated_by, v.deleted`

func scanVehicle(row pgx.Row) (*fleet.Vehicle, error) {
	var v fleet.Vehicle
	err := row.Scan(
		&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
		&v.Brand, &v.Model, &v.Year, &v.Colour, &v.FuelType, &v.TransmissionType,
		&v.CurrentMileage, &v.Status, &v.Designation, &v.AcquisitionDate,
		&v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
		&v.InsuranceExpiryDate, &v.RegistrationExpiryDate,
		&v.LastInspectionDate, &v.NextInspectionDueDate, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt, &v.CreatedBy, &v.UpdatedBy, &v.Deleted,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func scanVehicleRow(rows pgx.Rows) (*fleet.Vehicle, error) {
	var v fleet.Vehicle
	err := rows.Scan(
		&v.ID, &v.CarGroupID, &v.BranchID, &v.VIN, &v.LicencePlate,
		&v.Brand, &v.Model, &v.Year, &v.Colour, &v.FuelType, &v.TransmissionType,
		&v.CurrentMileage, &v.Status, &v.Designation, &v.AcquisitionDate,
		&v.OwnershipType, &v.LeaseDetails, &v.InsurancePolicyNumber,
		&v.InsuranceExpiryDate, &v.RegistrationExpiryDate,
		&v.LastInspectionDate, &v.NextInspectionDueDate, &v.Notes,
		&v.CreatedAt, &v.UpdatedAt, &v.DeletedAt, &v.CreatedBy, &v.UpdatedBy, &v.Deleted,
	)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// setExpiryWarning sets ExpiryWarning to true if insurance or registration expires within 30 days.
func setExpiryWarning(v *fleet.Vehicle) {
	threshold := time.Now().AddDate(0, 0, expiryWarningDays)
	if v.InsuranceExpiryDate != nil && !v.InsuranceExpiryDate.After(threshold) {
		v.ExpiryWarning = true
		return
	}
	if v.RegistrationExpiryDate != nil && !v.RegistrationExpiryDate.After(threshold) {
		v.ExpiryWarning = true
	}
}

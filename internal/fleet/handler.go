package fleet

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for the fleet module.
type Handler struct {
	carGroupRepo CarGroupRepository
	vehicleRepo  VehicleRepository
}

// NewHandler creates a new fleet HTTP handler.
func NewHandler(carGroupRepo CarGroupRepository, vehicleRepo VehicleRepository) *Handler {
	return &Handler{
		carGroupRepo: carGroupRepo,
		vehicleRepo:  vehicleRepo,
	}
}

// RegisterRoutes registers all fleet routes to the provided router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/car-groups", func(r chi.Router) {
		r.Post("/", h.CreateCarGroup)
		r.Get("/", h.ListCarGroups)
		r.Get("/{id}", h.GetCarGroup)
		r.Put("/{id}", h.UpdateCarGroup)
		r.Delete("/{id}", h.DeleteCarGroup)
	})

	r.Route("/vehicles", func(r chi.Router) {
		r.Post("/", h.CreateVehicle)
		r.Get("/", h.ListVehicles)
		r.Get("/{id}", h.GetVehicle)
		r.Put("/{id}", h.UpdateVehicle)
		r.Delete("/{id}", h.DeleteVehicle)
	})
}

// createCarGroupRequest is the request payload for creating a car group.
type createCarGroupRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
	CreatedBy    string  `json:"created_by"`
}

// updateCarGroupRequest is the request payload for updating a car group.
type updateCarGroupRequest struct {
	Name         string  `json:"name"`
	Description  *string `json:"description,omitempty"`
	SizeCategory *string `json:"size_category,omitempty"`
	UpdatedBy    string  `json:"updated_by"`
}

// deleteRequest carries the actor performing the soft-delete.
type deleteRequest struct {
	DeletedBy string `json:"deleted_by"`
}

// createVehicleRequest is the request payload for creating a vehicle.
type createVehicleRequest struct {
	CarGroupID             uuid.UUID          `json:"car_group_id"`
	BranchID               uuid.UUID          `json:"branch_id"`
	VIN                    string             `json:"vin"`
	LicencePlate           string             `json:"licence_plate"`
	Brand                  string             `json:"brand"`
	Model                  string             `json:"model"`
	Year                   int                `json:"year"`
	Colour                 *string            `json:"colour,omitempty"`
	FuelType               FuelType           `json:"fuel_type"`
	TransmissionType       TransmissionType   `json:"transmission_type"`
	CurrentMileage         int                `json:"current_mileage"`
	Status                 VehicleStatus      `json:"status"`
	Designation            VehicleDesignation `json:"designation"`
	AcquisitionDate        time.Time          `json:"acquisition_date"`
	OwnershipType          OwnershipType      `json:"ownership_type"`
	LeaseDetails           *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time         `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                  *string            `json:"notes,omitempty"`
	CreatedBy              string             `json:"created_by"`
}

// updateVehicleRequest is the request payload for updating a vehicle.
type updateVehicleRequest struct {
	CarGroupID             uuid.UUID          `json:"car_group_id"`
	BranchID               uuid.UUID          `json:"branch_id"`
	VIN                    string             `json:"vin"`
	LicencePlate           string             `json:"licence_plate"`
	Brand                  string             `json:"brand"`
	Model                  string             `json:"model"`
	Year                   int                `json:"year"`
	Colour                 *string            `json:"colour,omitempty"`
	FuelType               FuelType           `json:"fuel_type"`
	TransmissionType       TransmissionType   `json:"transmission_type"`
	CurrentMileage         int                `json:"current_mileage"`
	Status                 VehicleStatus      `json:"status"`
	Designation            VehicleDesignation `json:"designation"`
	AcquisitionDate        time.Time          `json:"acquisition_date"`
	OwnershipType          OwnershipType      `json:"ownership_type"`
	LeaseDetails           *string            `json:"lease_details,omitempty"`
	InsurancePolicyNumber  *string            `json:"insurance_policy_number,omitempty"`
	InsuranceExpiryDate    *time.Time         `json:"insurance_expiry_date,omitempty"`
	RegistrationExpiryDate *time.Time         `json:"registration_expiry_date,omitempty"`
	LastInspectionDate     *time.Time         `json:"last_inspection_date,omitempty"`
	NextInspectionDueDate  *time.Time         `json:"next_inspection_due_date,omitempty"`
	Notes                  *string            `json:"notes,omitempty"`
	UpdatedBy              string             `json:"updated_by"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// CreateCarGroup handles POST /car-groups.
func (h *Handler) CreateCarGroup(w http.ResponseWriter, r *http.Request) {
	var req createCarGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.CreatedBy == "" {
		writeError(w, http.StatusBadRequest, "name and created_by are required")
		return
	}

	now := time.Now().UTC()
	group := &CarGroup{
		ID:           uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		SizeCategory: req.SizeCategory,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    req.CreatedBy,
		UpdatedBy:    req.CreatedBy,
		Deleted:      false,
	}

	if err := h.carGroupRepo.Create(r.Context(), group); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create car group")
		return
	}
	writeJSON(w, http.StatusCreated, group)
}

// ListCarGroups handles GET /car-groups.
func (h *Handler) ListCarGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.carGroupRepo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list car groups")
		return
	}
	if groups == nil {
		groups = []*CarGroup{}
	}
	writeJSON(w, http.StatusOK, groups)
}

// GetCarGroup handles GET /car-groups/{id}.
func (h *Handler) GetCarGroup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	group, err := h.carGroupRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "car group not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get car group")
		return
	}
	writeJSON(w, http.StatusOK, group)
}

// UpdateCarGroup handles PUT /car-groups/{id}.
func (h *Handler) UpdateCarGroup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateCarGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.UpdatedBy == "" {
		writeError(w, http.StatusBadRequest, "name and updated_by are required")
		return
	}

	group := &CarGroup{
		ID:           id,
		Name:         req.Name,
		Description:  req.Description,
		SizeCategory: req.SizeCategory,
		UpdatedAt:    time.Now().UTC(),
		UpdatedBy:    req.UpdatedBy,
	}

	if err := h.carGroupRepo.Update(r.Context(), group); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "car group not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update car group")
		return
	}
	writeJSON(w, http.StatusOK, group)
}

// DeleteCarGroup handles DELETE /car-groups/{id}.
func (h *Handler) DeleteCarGroup(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req deleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DeletedBy == "" {
		writeError(w, http.StatusBadRequest, "deleted_by is required")
		return
	}

	if err := h.carGroupRepo.Delete(r.Context(), id, req.DeletedBy); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "car group not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete car group")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateVehicle handles POST /vehicles.
func (h *Handler) CreateVehicle(w http.ResponseWriter, r *http.Request) {
	var req createVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.VIN == "" || req.LicencePlate == "" || req.Brand == "" || req.Model == "" ||
		req.FuelType == "" || req.TransmissionType == "" || req.Status == "" ||
		req.Designation == "" || req.OwnershipType == "" || req.CreatedBy == "" {
		writeError(w, http.StatusBadRequest, "required fields are missing")
		return
	}

	now := time.Now().UTC()
	vehicle := &Vehicle{
		ID:                     uuid.New(),
		CarGroupID:             req.CarGroupID,
		BranchID:               req.BranchID,
		VIN:                    req.VIN,
		LicencePlate:           req.LicencePlate,
		Brand:                  req.Brand,
		Model:                  req.Model,
		Year:                   req.Year,
		Colour:                 req.Colour,
		FuelType:               req.FuelType,
		TransmissionType:       req.TransmissionType,
		CurrentMileage:         req.CurrentMileage,
		Status:                 req.Status,
		Designation:            req.Designation,
		AcquisitionDate:        req.AcquisitionDate,
		OwnershipType:          req.OwnershipType,
		LeaseDetails:           req.LeaseDetails,
		InsurancePolicyNumber:  req.InsurancePolicyNumber,
		InsuranceExpiryDate:    req.InsuranceExpiryDate,
		RegistrationExpiryDate: req.RegistrationExpiryDate,
		LastInspectionDate:     req.LastInspectionDate,
		NextInspectionDueDate:  req.NextInspectionDueDate,
		Notes:                  req.Notes,
		CreatedAt:              now,
		UpdatedAt:              now,
		CreatedBy:              req.CreatedBy,
		UpdatedBy:              req.CreatedBy,
		Deleted:                false,
	}

	if err := h.vehicleRepo.Create(r.Context(), vehicle); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create vehicle")
		return
	}
	writeJSON(w, http.StatusCreated, vehicle)
}

// ListVehicles handles GET /vehicles.
func (h *Handler) ListVehicles(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.vehicleRepo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}
	if vehicles == nil {
		vehicles = []*Vehicle{}
	}
	writeJSON(w, http.StatusOK, vehicles)
}

// GetVehicle handles GET /vehicles/{id}.
func (h *Handler) GetVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	vehicle, err := h.vehicleRepo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "vehicle not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get vehicle")
		return
	}
	writeJSON(w, http.StatusOK, vehicle)
}

// UpdateVehicle handles PUT /vehicles/{id}.
func (h *Handler) UpdateVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.VIN == "" || req.LicencePlate == "" || req.Brand == "" || req.Model == "" ||
		req.FuelType == "" || req.TransmissionType == "" || req.Status == "" ||
		req.Designation == "" || req.OwnershipType == "" || req.UpdatedBy == "" {
		writeError(w, http.StatusBadRequest, "required fields are missing")
		return
	}

	vehicle := &Vehicle{
		ID:                     id,
		CarGroupID:             req.CarGroupID,
		BranchID:               req.BranchID,
		VIN:                    req.VIN,
		LicencePlate:           req.LicencePlate,
		Brand:                  req.Brand,
		Model:                  req.Model,
		Year:                   req.Year,
		Colour:                 req.Colour,
		FuelType:               req.FuelType,
		TransmissionType:       req.TransmissionType,
		CurrentMileage:         req.CurrentMileage,
		Status:                 req.Status,
		Designation:            req.Designation,
		AcquisitionDate:        req.AcquisitionDate,
		OwnershipType:          req.OwnershipType,
		LeaseDetails:           req.LeaseDetails,
		InsurancePolicyNumber:  req.InsurancePolicyNumber,
		InsuranceExpiryDate:    req.InsuranceExpiryDate,
		RegistrationExpiryDate: req.RegistrationExpiryDate,
		LastInspectionDate:     req.LastInspectionDate,
		NextInspectionDueDate:  req.NextInspectionDueDate,
		Notes:                  req.Notes,
		UpdatedAt:              time.Now().UTC(),
		UpdatedBy:              req.UpdatedBy,
	}

	if err := h.vehicleRepo.Update(r.Context(), vehicle); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "vehicle not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update vehicle")
		return
	}
	writeJSON(w, http.StatusOK, vehicle)
}

// DeleteVehicle handles DELETE /vehicles/{id}.
func (h *Handler) DeleteVehicle(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req deleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.DeletedBy == "" {
		writeError(w, http.StatusBadRequest, "deleted_by is required")
		return
	}

	if err := h.vehicleRepo.Delete(r.Context(), id, req.DeletedBy); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "vehicle not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to delete vehicle")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

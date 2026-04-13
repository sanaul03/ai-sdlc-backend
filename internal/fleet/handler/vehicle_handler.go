package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// VehicleServicer defines the operations required by the vehicle handler.
type VehicleServicer interface {
	Create(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error)
	List(ctx context.Context, filter fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error)
	GetByID(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error)
	Update(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error)
	UpdateDesignation(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

// VehicleHandler handles HTTP requests for vehicle endpoints.
type VehicleHandler struct {
	svc VehicleServicer
}

// NewVehicleHandler constructs a VehicleHandler.
func NewVehicleHandler(svc VehicleServicer) *VehicleHandler {
	return &VehicleHandler{svc: svc}
}

// RegisterRoutes wires the handler into a ServeMux under /api/v1/vehicles.
func (h *VehicleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/vehicles", h.create)
	mux.HandleFunc("GET /api/v1/vehicles", h.list)
	mux.HandleFunc("GET /api/v1/vehicles/{id}", h.getByID)
	mux.HandleFunc("PUT /api/v1/vehicles/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/vehicles/{id}", h.delete)
	mux.HandleFunc("PATCH /api/v1/vehicles/{id}/designation", h.updateDesignation)
}

func (h *VehicleHandler) create(w http.ResponseWriter, r *http.Request) {
	var body createVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input, err := body.toInput(callerFromRequest(r))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	v, err := h.svc.Create(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, v)
}

func (h *VehicleHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := fleet.ListVehiclesFilter{}

	if v := q.Get("car_group_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid car_group_id")
			return
		}
		filter.CarGroupID = &id
	}
	if v := q.Get("branch_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid branch_id")
			return
		}
		filter.BranchID = &id
	}
	if v := q.Get("status"); v != "" {
		s := fleet.VehicleStatus(v)
		filter.Status = &s
	}
	if v := q.Get("designation"); v != "" {
		d := fleet.VehicleDesignation(v)
		filter.Designation = &d
	}
	if v := q.Get("fuel_type"); v != "" {
		ft := fleet.FuelType(v)
		filter.FuelType = &ft
	}
	if v := q.Get("transmission_type"); v != "" {
		tt := fleet.TransmissionType(v)
		filter.TransmissionType = &tt
	}
	if v := q.Get("expiry_warning"); v != "" {
		b := v == "true"
		filter.ExpiryWarning = &b
	}
	if v := q.Get("page"); v != "" {
		p, _ := strconv.Atoi(v)
		filter.Page = p
	}
	if v := q.Get("page_size"); v != "" {
		ps, _ := strconv.Atoi(v)
		filter.PageSize = ps
	}

	vehicles, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	if vehicles == nil {
		vehicles = []*fleet.Vehicle{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"data":  vehicles,
		"total": total,
	})
}

func (h *VehicleHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	v, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *VehicleHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body updateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := body.toInput(callerFromRequest(r))
	v, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

func (h *VehicleHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id, callerFromRequest(r)); err != nil {
		handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *VehicleHandler) updateDesignation(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Designation string `json:"designation"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := fleet.UpdateDesignationInput{
		Designation: fleet.VehicleDesignation(body.Designation),
		UpdatedBy:   callerFromRequest(r),
	}

	v, err := h.svc.UpdateDesignation(r.Context(), id, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, v)
}

// --- request DTOs ---

type createVehicleRequest struct {
	CarGroupID              string  `json:"car_group_id"`
	BranchID                string  `json:"branch_id"`
	VIN                     string  `json:"vin"`
	LicencePlate            string  `json:"licence_plate"`
	Brand                   string  `json:"brand"`
	Model                   string  `json:"model"`
	Year                    int     `json:"year"`
	Colour                  *string `json:"colour"`
	FuelType                string  `json:"fuel_type"`
	TransmissionType        string  `json:"transmission_type"`
	CurrentMileage          int     `json:"current_mileage"`
	Status                  string  `json:"status"`
	Designation             string  `json:"designation"`
	AcquisitionDate         string  `json:"acquisition_date"`
	OwnershipType           string  `json:"ownership_type"`
	LeaseDetails            *string `json:"lease_details"`
	InsurancePolicyNumber   *string `json:"insurance_policy_number"`
	InsuranceExpiryDate     *string `json:"insurance_expiry_date"`
	RegistrationExpiryDate  *string `json:"registration_expiry_date"`
	LastInspectionDate      *string `json:"last_inspection_date"`
	NextInspectionDueDate   *string `json:"next_inspection_due_date"`
	Notes                   *string `json:"notes"`
}

func (req createVehicleRequest) toInput(createdBy string) (fleet.CreateVehicleInput, error) {
	carGroupID, err := uuid.Parse(req.CarGroupID)
	if err != nil {
		return fleet.CreateVehicleInput{}, &parseError{"car_group_id", err}
	}
	branchID, err := uuid.Parse(req.BranchID)
	if err != nil {
		return fleet.CreateVehicleInput{}, &parseError{"branch_id", err}
	}
	acqDate, err := time.Parse("2006-01-02", req.AcquisitionDate)
	if err != nil {
		return fleet.CreateVehicleInput{}, &parseError{"acquisition_date", err}
	}

	input := fleet.CreateVehicleInput{
		CarGroupID:            carGroupID,
		BranchID:              branchID,
		VIN:                   req.VIN,
		LicencePlate:          req.LicencePlate,
		Brand:                 req.Brand,
		Model:                 req.Model,
		Year:                  req.Year,
		Colour:                req.Colour,
		FuelType:              fleet.FuelType(req.FuelType),
		TransmissionType:      fleet.TransmissionType(req.TransmissionType),
		CurrentMileage:        req.CurrentMileage,
		Status:                fleet.VehicleStatus(req.Status),
		Designation:           fleet.VehicleDesignation(req.Designation),
		AcquisitionDate:       acqDate,
		OwnershipType:         req.OwnershipType,
		LeaseDetails:          req.LeaseDetails,
		InsurancePolicyNumber: req.InsurancePolicyNumber,
		Notes:                 req.Notes,
		CreatedBy:             createdBy,
	}

	if req.InsuranceExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.InsuranceExpiryDate)
		if err != nil {
			return fleet.CreateVehicleInput{}, &parseError{"insurance_expiry_date", err}
		}
		input.InsuranceExpiryDate = &t
	}
	if req.RegistrationExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.RegistrationExpiryDate)
		if err != nil {
			return fleet.CreateVehicleInput{}, &parseError{"registration_expiry_date", err}
		}
		input.RegistrationExpiryDate = &t
	}
	if req.LastInspectionDate != nil {
		t, err := time.Parse("2006-01-02", *req.LastInspectionDate)
		if err != nil {
			return fleet.CreateVehicleInput{}, &parseError{"last_inspection_date", err}
		}
		input.LastInspectionDate = &t
	}
	if req.NextInspectionDueDate != nil {
		t, err := time.Parse("2006-01-02", *req.NextInspectionDueDate)
		if err != nil {
			return fleet.CreateVehicleInput{}, &parseError{"next_inspection_due_date", err}
		}
		input.NextInspectionDueDate = &t
	}

	return input, nil
}

type updateVehicleRequest struct {
	CarGroupID             *string `json:"car_group_id"`
	BranchID               *string `json:"branch_id"`
	VIN                    *string `json:"vin"`
	LicencePlate           *string `json:"licence_plate"`
	Brand                  *string `json:"brand"`
	Model                  *string `json:"model"`
	Year                   *int    `json:"year"`
	Colour                 *string `json:"colour"`
	FuelType               *string `json:"fuel_type"`
	TransmissionType       *string `json:"transmission_type"`
	CurrentMileage         *int    `json:"current_mileage"`
	Designation            *string `json:"designation"`
	AcquisitionDate        *string `json:"acquisition_date"`
	OwnershipType          *string `json:"ownership_type"`
	LeaseDetails           *string `json:"lease_details"`
	InsurancePolicyNumber  *string `json:"insurance_policy_number"`
	InsuranceExpiryDate    *string `json:"insurance_expiry_date"`
	RegistrationExpiryDate *string `json:"registration_expiry_date"`
	LastInspectionDate     *string `json:"last_inspection_date"`
	NextInspectionDueDate  *string `json:"next_inspection_due_date"`
	Notes                  *string `json:"notes"`
}

func (req updateVehicleRequest) toInput(updatedBy string) fleet.UpdateVehicleInput {
	input := fleet.UpdateVehicleInput{UpdatedBy: updatedBy}

	if req.CarGroupID != nil {
		id, err := uuid.Parse(*req.CarGroupID)
		if err == nil {
			input.CarGroupID = &id
		}
	}
	if req.BranchID != nil {
		id, err := uuid.Parse(*req.BranchID)
		if err == nil {
			input.BranchID = &id
		}
	}
	input.VIN = req.VIN
	input.LicencePlate = req.LicencePlate
	input.Brand = req.Brand
	input.Model = req.Model
	input.Year = req.Year
	input.Colour = req.Colour
	input.CurrentMileage = req.CurrentMileage
	input.OwnershipType = req.OwnershipType
	input.LeaseDetails = req.LeaseDetails
	input.InsurancePolicyNumber = req.InsurancePolicyNumber
	input.Notes = req.Notes

	if req.FuelType != nil {
		ft := fleet.FuelType(*req.FuelType)
		input.FuelType = &ft
	}
	if req.TransmissionType != nil {
		tt := fleet.TransmissionType(*req.TransmissionType)
		input.TransmissionType = &tt
	}
	if req.Designation != nil {
		d := fleet.VehicleDesignation(*req.Designation)
		input.Designation = &d
	}
	if req.AcquisitionDate != nil {
		t, err := time.Parse("2006-01-02", *req.AcquisitionDate)
		if err == nil {
			input.AcquisitionDate = &t
		}
	}
	if req.InsuranceExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.InsuranceExpiryDate)
		if err == nil {
			input.InsuranceExpiryDate = &t
		}
	}
	if req.RegistrationExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.RegistrationExpiryDate)
		if err == nil {
			input.RegistrationExpiryDate = &t
		}
	}
	if req.LastInspectionDate != nil {
		t, err := time.Parse("2006-01-02", *req.LastInspectionDate)
		if err == nil {
			input.LastInspectionDate = &t
		}
	}
	if req.NextInspectionDueDate != nil {
		t, err := time.Parse("2006-01-02", *req.NextInspectionDueDate)
		if err == nil {
			input.NextInspectionDueDate = &t
		}
	}
	return input
}

// parseError describes a request field that could not be parsed.
type parseError struct {
	field string
	err   error
}

func (e *parseError) Error() string {
	return "invalid value for field '" + e.field + "': " + e.err.Error()
}

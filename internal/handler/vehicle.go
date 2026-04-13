package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"
)

// VehicleHandler handles HTTP requests for vehicle resources.
type VehicleHandler struct {
	repo repository.Vehicle
}

// NewVehicleHandler creates a new VehicleHandler.
func NewVehicleHandler(repo repository.Vehicle) *VehicleHandler {
	return &VehicleHandler{repo: repo}
}

// Create handles POST /vehicles.
func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CompanyID == 0 || req.VehicleTypeID == 0 || req.RegistrationNumber == "" {
		respondError(w, http.StatusBadRequest, "company_id, vehicle_type_id and registration_number are required")
		return
	}

	vehicle, err := h.repo.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create vehicle")
		return
	}
	respondJSON(w, http.StatusCreated, vehicle)
}

// GetByID handles GET /vehicles/{id}.
func (h *VehicleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	vehicle, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "vehicle not found")
		return
	}
	respondJSON(w, http.StatusOK, vehicle)
}

// ListByCompany handles GET /companies/{id}/vehicles.
func (h *VehicleHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyID, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid company id")
		return
	}

	vehicles, err := h.repo.ListByCompany(r.Context(), companyID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list vehicles")
		return
	}
	respondJSON(w, http.StatusOK, vehicles)
}

// Update handles PUT /vehicles/{id}.
func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vehicle, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update vehicle")
		return
	}
	respondJSON(w, http.StatusOK, vehicle)
}

// Delete handles DELETE /vehicles/{id}.
func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete vehicle")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

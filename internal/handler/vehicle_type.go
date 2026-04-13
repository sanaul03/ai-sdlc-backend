package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"
)

// VehicleTypeHandler handles HTTP requests for vehicle type resources.
type VehicleTypeHandler struct {
	repo repository.VehicleType
}

// NewVehicleTypeHandler creates a new VehicleTypeHandler.
func NewVehicleTypeHandler(repo repository.VehicleType) *VehicleTypeHandler {
	return &VehicleTypeHandler{repo: repo}
}

// Create handles POST /vehicle-types.
func (h *VehicleTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVehicleTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CategoryID == 0 || req.Name == "" || req.Code == "" {
		respondError(w, http.StatusBadRequest, "category_id, name and code are required")
		return
	}

	vt, err := h.repo.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create vehicle type")
		return
	}
	respondJSON(w, http.StatusCreated, vt)
}

// GetByID handles GET /vehicle-types/{id}.
func (h *VehicleTypeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	vt, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "vehicle type not found")
		return
	}
	respondJSON(w, http.StatusOK, vt)
}

// ListByCategory handles GET /vehicle-categories/{id}/vehicle-types.
func (h *VehicleTypeHandler) ListByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	types, err := h.repo.ListByCategory(r.Context(), categoryID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list vehicle types")
		return
	}
	respondJSON(w, http.StatusOK, types)
}

// Update handles PUT /vehicle-types/{id}.
func (h *VehicleTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateVehicleTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vt, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update vehicle type")
		return
	}
	respondJSON(w, http.StatusOK, vt)
}

// Delete handles DELETE /vehicle-types/{id}.
func (h *VehicleTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete vehicle type")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

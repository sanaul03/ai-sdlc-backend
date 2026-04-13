package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"
)

// VehicleCategoryHandler handles HTTP requests for vehicle category resources.
type VehicleCategoryHandler struct {
	repo repository.VehicleCategory
}

// NewVehicleCategoryHandler creates a new VehicleCategoryHandler.
func NewVehicleCategoryHandler(repo repository.VehicleCategory) *VehicleCategoryHandler {
	return &VehicleCategoryHandler{repo: repo}
}

// Create handles POST /vehicle-categories.
func (h *VehicleCategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateVehicleCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Code == "" {
		respondError(w, http.StatusBadRequest, "name and code are required")
		return
	}

	vc, err := h.repo.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create vehicle category")
		return
	}
	respondJSON(w, http.StatusCreated, vc)
}

// GetByID handles GET /vehicle-categories/{id}.
func (h *VehicleCategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	vc, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "vehicle category not found")
		return
	}
	respondJSON(w, http.StatusOK, vc)
}

// List handles GET /vehicle-categories.
func (h *VehicleCategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list vehicle categories")
		return
	}
	respondJSON(w, http.StatusOK, categories)
}

// Update handles PUT /vehicle-categories/{id}.
func (h *VehicleCategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateVehicleCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vc, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update vehicle category")
		return
	}
	respondJSON(w, http.StatusOK, vc)
}

// Delete handles DELETE /vehicle-categories/{id}.
func (h *VehicleCategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete vehicle category")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

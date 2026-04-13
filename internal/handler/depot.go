package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"
)

// DepotHandler handles HTTP requests for depot resources.
type DepotHandler struct {
	repo repository.Depot
}

// NewDepotHandler creates a new DepotHandler.
func NewDepotHandler(repo repository.Depot) *DepotHandler {
	return &DepotHandler{repo: repo}
}

// Create handles POST /depots.
func (h *DepotHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateDepotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.CompanyID == 0 || req.Name == "" || req.Code == "" {
		respondError(w, http.StatusBadRequest, "company_id, name and code are required")
		return
	}

	depot, err := h.repo.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create depot")
		return
	}
	respondJSON(w, http.StatusCreated, depot)
}

// GetByID handles GET /depots/{id}.
func (h *DepotHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	depot, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "depot not found")
		return
	}
	respondJSON(w, http.StatusOK, depot)
}

// ListByCompany handles GET /companies/{id}/depots.
func (h *DepotHandler) ListByCompany(w http.ResponseWriter, r *http.Request) {
	companyID, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid company id")
		return
	}

	depots, err := h.repo.ListByCompany(r.Context(), companyID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list depots")
		return
	}
	respondJSON(w, http.StatusOK, depots)
}

// Update handles PUT /depots/{id}.
func (h *DepotHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateDepotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	depot, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update depot")
		return
	}
	respondJSON(w, http.StatusOK, depot)
}

// Delete handles DELETE /depots/{id}.
func (h *DepotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete depot")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

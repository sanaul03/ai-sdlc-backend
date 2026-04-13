package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sanaul03/ai-sdlc-backend/internal/model"
	"github.com/sanaul03/ai-sdlc-backend/internal/repository"
)

// CompanyHandler handles HTTP requests for company resources.
type CompanyHandler struct {
	repo repository.Company
}

// NewCompanyHandler creates a new CompanyHandler.
func NewCompanyHandler(repo repository.Company) *CompanyHandler {
	return &CompanyHandler{repo: repo}
}

// Create handles POST /companies.
func (h *CompanyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Code == "" {
		respondError(w, http.StatusBadRequest, "name and code are required")
		return
	}

	company, err := h.repo.Create(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create company")
		return
	}
	respondJSON(w, http.StatusCreated, company)
}

// GetByID handles GET /companies/{id}.
func (h *CompanyHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	company, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "company not found")
		return
	}
	respondJSON(w, http.StatusOK, company)
}

// List handles GET /companies.
func (h *CompanyHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list companies")
		return
	}
	respondJSON(w, http.StatusOK, companies)
}

// Update handles PUT /companies/{id}.
func (h *CompanyHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req model.UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	company, err := h.repo.Update(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to update company")
		return
	}
	respondJSON(w, http.StatusOK, company)
}

// Delete handles DELETE /companies/{id}.
func (h *CompanyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathIDParam(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to delete company")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// pathIDParam extracts the "id" path value from the request (Go 1.22+ ServeMux).
func pathIDParam(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

// respondJSON writes a JSON response with the given status code.
func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// respondError writes a JSON error response.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

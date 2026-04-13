package vehicle

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/middleware"
)

// Handler exposes HTTP endpoints for the vehicle resource.
type Handler struct {
	svc Service
}

// NewHandler creates a Handler with the provided Service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts all vehicle endpoints under the given chi.Router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	h.RegisterWriteRoutes(r)
	h.RegisterReadRoutes(r)
}

// RegisterWriteRoutes mounts the mutating endpoints (POST, PUT, DELETE, PATCH).
func (h *Handler) RegisterWriteRoutes(r chi.Router) {
	r.Post("/vehicles", h.create)
	r.Put("/vehicles/{id}", h.update)
	r.Delete("/vehicles/{id}", h.delete)
	r.Patch("/vehicles/{id}/designation", h.updateDesignation)
}

// RegisterReadRoutes mounts the read-only endpoints (GET).
func (h *Handler) RegisterReadRoutes(r chi.Router) {
	r.Get("/vehicles", h.list)
	r.Get("/vehicles/{id}", h.get)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v, err := h.svc.Create(r.Context(), req, claims.Sub)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, v)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := ListFilter{
		Status:           q.Get("status"),
		Designation:      q.Get("designation"),
		FuelType:         q.Get("fuel_type"),
		TransmissionType: q.Get("transmission_type"),
	}

	if v := q.Get("car_group_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			filter.CarGroupID = &id
		}
	}
	if v := q.Get("branch_id"); v != "" {
		id, err := uuid.Parse(v)
		if err == nil {
			filter.BranchID = &id
		}
	}
	if v := q.Get("expiry_warning"); v != "" {
		filter.ExpiryWarning, _ = strconv.ParseBool(v)
	}
	if v := q.Get("page"); v != "" {
		filter.Page, _ = strconv.Atoi(v)
	}
	if v := q.Get("page_size"); v != "" {
		filter.PageSize, _ = strconv.Atoi(v)
	}

	page, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, page)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
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

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v, err := h.svc.Update(r.Context(), id, req, claims.Sub)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, v)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id, claims.Sub); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) updateDesignation(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req DesignationUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	v, err := h.svc.UpdateDesignation(r.Context(), id, req, claims.Sub)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, v)
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

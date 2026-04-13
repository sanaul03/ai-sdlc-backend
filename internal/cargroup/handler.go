package cargroup

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/middleware"
)

// Handler exposes HTTP endpoints for the car group resource.
type Handler struct {
	svc Service
}

// NewHandler creates a Handler with the provided Service.
func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes mounts all car group endpoints under the given chi.Router.
// The caller is responsible for applying authentication middleware beforehand.
func (h *Handler) RegisterRoutes(r chi.Router) {
	h.RegisterWriteRoutes(r)
	h.RegisterReadRoutes(r)
}

// RegisterWriteRoutes mounts the mutating endpoints (POST, PUT, DELETE).
func (h *Handler) RegisterWriteRoutes(r chi.Router) {
	r.Post("/car-groups", h.create)
	r.Put("/car-groups/{id}", h.update)
	r.Delete("/car-groups/{id}", h.delete)
}

// RegisterReadRoutes mounts the read-only endpoints (GET).
func (h *Handler) RegisterReadRoutes(r chi.Router) {
	r.Get("/car-groups", h.list)
	r.Get("/car-groups/{id}", h.get)
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

	cg, err := h.svc.Create(r.Context(), req, claims.Sub)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, cg)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")

	deletedStr := r.URL.Query().Get("deleted")
	deleted, _ := strconv.ParseBool(deletedStr)

	filter := ListFilter{Q: q, Deleted: deleted}

	groups, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	if groups == nil {
		groups = []CarGroup{}
	}
	writeJSON(w, http.StatusOK, groups)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	cg, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, cg)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	cg, err := h.svc.Update(r.Context(), id, req, claims.Sub)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, cg)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing authentication")
		return
	}

	id, err := parseUUID(chi.URLParam(r, "id"))
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

// --- helpers ---

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

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
	case errors.Is(err, ErrHasVehicles):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

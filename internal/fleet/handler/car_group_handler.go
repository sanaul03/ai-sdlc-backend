package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// CarGroupServicer defines the operations required by the car group handler.
type CarGroupServicer interface {
	Create(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error)
	List(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error)
	GetByID(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error)
	Update(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error)
	Delete(ctx context.Context, id uuid.UUID, deletedBy string) error
}

// CarGroupHandler handles HTTP requests for car group endpoints.
type CarGroupHandler struct {
	svc CarGroupServicer
}

// NewCarGroupHandler constructs a CarGroupHandler.
func NewCarGroupHandler(svc CarGroupServicer) *CarGroupHandler {
	return &CarGroupHandler{svc: svc}
}

// RegisterRoutes wires the handler into a ServeMux under /api/v1/car-groups.
func (h *CarGroupHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/car-groups", h.create)
	mux.HandleFunc("GET /api/v1/car-groups", h.list)
	mux.HandleFunc("GET /api/v1/car-groups/{id}", h.getByID)
	mux.HandleFunc("PUT /api/v1/car-groups/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/car-groups/{id}", h.delete)
}

func (h *CarGroupHandler) create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name         string  `json:"name"`
		Description  *string `json:"description"`
		SizeCategory *string `json:"size_category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	caller := callerFromRequest(r)
	input := fleet.CreateCarGroupInput{
		Name:         body.Name,
		Description:  body.Description,
		SizeCategory: body.SizeCategory,
		CreatedBy:    caller,
	}

	g, err := h.svc.Create(r.Context(), input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, g)
}

func (h *CarGroupHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	deletedParam := r.URL.Query().Get("deleted")

	filter := fleet.ListCarGroupsFilter{}
	if q != "" {
		filter.Q = &q
	}
	filter.IncludeDeleted = deletedParam == "true"

	groups, err := h.svc.List(r.Context(), filter)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	if groups == nil {
		groups = []*fleet.CarGroup{}
	}
	writeJSON(w, http.StatusOK, groups)
}

func (h *CarGroupHandler) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	g, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *CarGroupHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUUID(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Name         *string `json:"name"`
		Description  *string `json:"description"`
		SizeCategory *string `json:"size_category"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input := fleet.UpdateCarGroupInput{
		Name:         body.Name,
		Description:  body.Description,
		SizeCategory: body.SizeCategory,
		UpdatedBy:    callerFromRequest(r),
	}

	g, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, g)
}

func (h *CarGroupHandler) delete(w http.ResponseWriter, r *http.Request) {
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

// --- shared helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, fleet.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, fleet.ErrValidation):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, fleet.ErrConflict):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// callerFromRequest extracts the caller identity from the request context.
// In production this would be populated from a JWT middleware; here we fall
// back to a header value for development convenience.
func callerFromRequest(r *http.Request) string {
	if caller := r.Header.Get("X-User-ID"); caller != "" {
		return caller
	}
	return "system"
}

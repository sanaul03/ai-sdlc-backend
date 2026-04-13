package vehicle_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/middleware"
	"github.com/sanaul03/ai-sdlc-backend/internal/vehicle"
)

// --- mock service ---

type mockService struct {
	createFn            func(ctx context.Context, req vehicle.CreateRequest, createdBy string) (vehicle.Vehicle, error)
	listFn              func(ctx context.Context, filter vehicle.ListFilter) (vehicle.Page, error)
	getByIDFn           func(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error)
	updateFn            func(ctx context.Context, id uuid.UUID, req vehicle.UpdateRequest, updatedBy string) (vehicle.Vehicle, error)
	updateDesignationFn func(ctx context.Context, id uuid.UUID, req vehicle.DesignationUpdateRequest, updatedBy string) (vehicle.Vehicle, error)
	deleteFn            func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockService) Create(ctx context.Context, req vehicle.CreateRequest, createdBy string) (vehicle.Vehicle, error) {
	return m.createFn(ctx, req, createdBy)
}
func (m *mockService) List(ctx context.Context, f vehicle.ListFilter) (vehicle.Page, error) {
	return m.listFn(ctx, f)
}
func (m *mockService) GetByID(ctx context.Context, id uuid.UUID) (vehicle.Vehicle, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockService) Update(ctx context.Context, id uuid.UUID, req vehicle.UpdateRequest, updatedBy string) (vehicle.Vehicle, error) {
	return m.updateFn(ctx, id, req, updatedBy)
}
func (m *mockService) UpdateDesignation(ctx context.Context, id uuid.UUID, req vehicle.DesignationUpdateRequest, updatedBy string) (vehicle.Vehicle, error) {
	return m.updateDesignationFn(ctx, id, req, updatedBy)
}
func (m *mockService) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.deleteFn(ctx, id, deletedBy)
}

// withClaims injects fake JWT claims into the request context.
func withClaims(r *http.Request) *http.Request {
	claims := &middleware.Claims{Sub: "test-user", Roles: []string{"FLEET_MANAGER"}}
	ctx := context.WithValue(r.Context(), middleware.ClaimsKey, claims)
	return r.WithContext(ctx)
}

// --- tests ---

func TestVehicleHandler_Create_201(t *testing.T) {
	v := newSampleVehicle()
	svc := &mockService{
		createFn: func(_ context.Context, _ vehicle.CreateRequest, _ string) (vehicle.Vehicle, error) {
			return v, nil
		},
	}
	h := vehicle.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	reqBody, _ := json.Marshal(newValidCreateRequest())
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewReader(reqBody))
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestVehicleHandler_Get_404(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (vehicle.Vehicle, error) {
			return vehicle.Vehicle{}, vehicle.ErrNotFound
		},
	}
	h := vehicle.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterReadRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/vehicles/"+uuid.New().String(), nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestVehicleHandler_UpdateDesignation_400(t *testing.T) {
	svc := &mockService{
		updateDesignationFn: func(_ context.Context, _ uuid.UUID, req vehicle.DesignationUpdateRequest, _ string) (vehicle.Vehicle, error) {
			return vehicle.Vehicle{}, vehicle.ErrInvalidInput
		},
	}
	h := vehicle.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	body, _ := json.Marshal(vehicle.DesignationUpdateRequest{Designation: "bad"})
	req := httptest.NewRequest(http.MethodPatch, "/vehicles/"+uuid.New().String()+"/designation", bytes.NewReader(body))
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestVehicleHandler_Delete_204(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error { return nil },
	}
	h := vehicle.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/vehicles/"+uuid.New().String(), nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestVehicleHandler_List_200(t *testing.T) {
	svc := &mockService{
		listFn: func(_ context.Context, _ vehicle.ListFilter) (vehicle.Page, error) {
			return vehicle.Page{Items: []vehicle.Vehicle{newSampleVehicle()}, Total: 1, Page: 1, PageSize: 20}, nil
		},
	}
	h := vehicle.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterReadRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/vehicles", nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

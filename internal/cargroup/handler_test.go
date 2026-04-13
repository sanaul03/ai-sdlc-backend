package cargroup_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/cargroup"
	"github.com/sanaul03/ai-sdlc-backend/internal/platform/middleware"
)

// --- mock service ---

type mockService struct {
	createFn  func(ctx context.Context, req cargroup.CreateRequest, createdBy string) (cargroup.CarGroup, error)
	listFn    func(ctx context.Context, filter cargroup.ListFilter) ([]cargroup.CarGroup, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (cargroup.CarGroup, error)
	updateFn  func(ctx context.Context, id uuid.UUID, req cargroup.UpdateRequest, updatedBy string) (cargroup.CarGroup, error)
	deleteFn  func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockService) Create(ctx context.Context, req cargroup.CreateRequest, createdBy string) (cargroup.CarGroup, error) {
	return m.createFn(ctx, req, createdBy)
}
func (m *mockService) List(ctx context.Context, filter cargroup.ListFilter) ([]cargroup.CarGroup, error) {
	return m.listFn(ctx, filter)
}
func (m *mockService) GetByID(ctx context.Context, id uuid.UUID) (cargroup.CarGroup, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockService) Update(ctx context.Context, id uuid.UUID, req cargroup.UpdateRequest, updatedBy string) (cargroup.CarGroup, error) {
	return m.updateFn(ctx, id, req, updatedBy)
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

func TestHandler_Create_201(t *testing.T) {
	sample := newSampleGroup()
	svc := &mockService{
		createFn: func(_ context.Context, _ cargroup.CreateRequest, _ string) (cargroup.CarGroup, error) {
			return sample, nil
		},
	}
	h := cargroup.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	body, _ := json.Marshal(cargroup.CreateRequest{Name: sample.Name})
	req := httptest.NewRequest(http.MethodPost, "/car-groups", bytes.NewReader(body))
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestHandler_Get_404(t *testing.T) {
	svc := &mockService{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (cargroup.CarGroup, error) {
			return cargroup.CarGroup{}, cargroup.ErrNotFound
		},
	}
	h := cargroup.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterReadRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/car-groups/"+uuid.New().String(), nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_Delete_Conflict(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return cargroup.ErrHasVehicles
		},
	}
	h := cargroup.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/car-groups/"+uuid.New().String(), nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestHandler_Delete_204(t *testing.T) {
	svc := &mockService{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	h := cargroup.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterWriteRoutes(r)

	req := httptest.NewRequest(http.MethodDelete, "/car-groups/"+uuid.New().String(), nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestHandler_List_200(t *testing.T) {
	svc := &mockService{
		listFn: func(_ context.Context, _ cargroup.ListFilter) ([]cargroup.CarGroup, error) {
			return []cargroup.CarGroup{newSampleGroup()}, nil
		},
	}
	h := cargroup.NewHandler(svc)

	r := chi.NewRouter()
	h.RegisterReadRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/car-groups", nil)
	req = withClaims(req)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

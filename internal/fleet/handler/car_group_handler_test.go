package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet/handler"
)

// --- mock service ---

type mockCarGroupSvc struct {
	createFn  func(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error)
	listFn    func(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error)
	getByIDFn func(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error)
	updateFn  func(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error)
	deleteFn  func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockCarGroupSvc) Create(ctx context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error) {
	return m.createFn(ctx, input)
}
func (m *mockCarGroupSvc) List(ctx context.Context, filter fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error) {
	return m.listFn(ctx, filter)
}
func (m *mockCarGroupSvc) GetByID(ctx context.Context, id uuid.UUID) (*fleet.CarGroup, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockCarGroupSvc) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateCarGroupInput) (*fleet.CarGroup, error) {
	return m.updateFn(ctx, id, input)
}
func (m *mockCarGroupSvc) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.deleteFn(ctx, id, deletedBy)
}

// --- helpers ---

func newTestMux(svc handler.CarGroupServicer) *http.ServeMux {
	mux := http.NewServeMux()
	h := handler.NewCarGroupHandler(svc)
	h.RegisterRoutes(mux)
	return mux
}

// --- tests ---

func TestCarGroupHandler_Create_201(t *testing.T) {
	id := uuid.New()
	svc := &mockCarGroupSvc{
		createFn: func(_ context.Context, input fleet.CreateCarGroupInput) (*fleet.CarGroup, error) {
			return &fleet.CarGroup{ID: id, Name: input.Name}, nil
		},
	}
	mux := newTestMux(svc)

	body, _ := json.Marshal(map[string]string{"name": "Economy Sedan"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/car-groups", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCarGroupHandler_Create_BadJSON(t *testing.T) {
	mux := newTestMux(&mockCarGroupSvc{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/car-groups", bytes.NewBufferString("{invalid"))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCarGroupHandler_List_200(t *testing.T) {
	svc := &mockCarGroupSvc{
		listFn: func(_ context.Context, _ fleet.ListCarGroupsFilter) ([]*fleet.CarGroup, error) {
			return []*fleet.CarGroup{{ID: uuid.New(), Name: "Luxury SUV"}}, nil
		},
	}
	mux := newTestMux(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/car-groups", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCarGroupHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockCarGroupSvc{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*fleet.CarGroup, error) {
			return nil, fleet.ErrNotFound
		},
	}
	mux := newTestMux(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/car-groups/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestCarGroupHandler_GetByID_InvalidUUID(t *testing.T) {
	mux := newTestMux(&mockCarGroupSvc{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/car-groups/not-a-uuid", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCarGroupHandler_Delete_NoContent(t *testing.T) {
	svc := &mockCarGroupSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	mux := newTestMux(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/car-groups/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
}

func TestCarGroupHandler_Delete_Conflict(t *testing.T) {
	svc := &mockCarGroupSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return fleet.ErrConflict
		},
	}
	mux := newTestMux(svc)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/car-groups/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", rr.Code)
	}
}

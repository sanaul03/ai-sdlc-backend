package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
	"github.com/sanaul03/ai-sdlc-backend/internal/fleet/handler"
)

// --- mock service ---

type mockVehicleSvc struct {
	createFn            func(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error)
	listFn              func(ctx context.Context, filter fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error)
	getByIDFn           func(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error)
	updateFn            func(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error)
	updateDesignationFn func(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error)
	deleteFn            func(ctx context.Context, id uuid.UUID, deletedBy string) error
}

func (m *mockVehicleSvc) Create(ctx context.Context, input fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
	return m.createFn(ctx, input)
}
func (m *mockVehicleSvc) List(ctx context.Context, f fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error) {
	return m.listFn(ctx, f)
}
func (m *mockVehicleSvc) GetByID(ctx context.Context, id uuid.UUID) (*fleet.Vehicle, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockVehicleSvc) Update(ctx context.Context, id uuid.UUID, input fleet.UpdateVehicleInput) (*fleet.Vehicle, error) {
	return m.updateFn(ctx, id, input)
}
func (m *mockVehicleSvc) UpdateDesignation(ctx context.Context, id uuid.UUID, input fleet.UpdateDesignationInput) (*fleet.Vehicle, error) {
	return m.updateDesignationFn(ctx, id, input)
}
func (m *mockVehicleSvc) Delete(ctx context.Context, id uuid.UUID, deletedBy string) error {
	return m.deleteFn(ctx, id, deletedBy)
}

// --- helpers ---

func newVehicleTestMux(svc handler.VehicleServicer) *http.ServeMux {
	mux := http.NewServeMux()
	h := handler.NewVehicleHandler(svc)
	h.RegisterRoutes(mux)
	return mux
}

func validCreateBody() map[string]any {
	return map[string]any{
		"car_group_id":      uuid.New().String(),
		"branch_id":         uuid.New().String(),
		"vin":               "1HGBH41JXMN109186",
		"licence_plate":     "ABC-1234",
		"brand":             "Toyota",
		"model":             "Corolla",
		"year":              2022,
		"fuel_type":         "petrol",
		"transmission_type": "automatic",
		"current_mileage":   0,
		"status":            "unavailable",
		"designation":       "rental_only",
		"acquisition_date":  time.Now().AddDate(-1, 0, 0).Format("2006-01-02"),
		"ownership_type":    "owned",
	}
}

// --- tests ---

func TestVehicleHandler_Create_201(t *testing.T) {
	id := uuid.New()
	svc := &mockVehicleSvc{
		createFn: func(_ context.Context, _ fleet.CreateVehicleInput) (*fleet.Vehicle, error) {
			return &fleet.Vehicle{ID: id}, nil
		},
	}
	mux := newVehicleTestMux(svc)

	body, _ := json.Marshal(validCreateBody())
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vehicles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestVehicleHandler_Create_BadJSON(t *testing.T) {
	mux := newVehicleTestMux(&mockVehicleSvc{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vehicles", bytes.NewBufferString("{bad"))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestVehicleHandler_Create_InvalidCarGroupID(t *testing.T) {
	mux := newVehicleTestMux(&mockVehicleSvc{})
	b := validCreateBody()
	b["car_group_id"] = "not-a-uuid"
	body, _ := json.Marshal(b)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/vehicles", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestVehicleHandler_List_200(t *testing.T) {
	svc := &mockVehicleSvc{
		listFn: func(_ context.Context, _ fleet.ListVehiclesFilter) ([]*fleet.Vehicle, int, error) {
			return []*fleet.Vehicle{{ID: uuid.New()}}, 1, nil
		},
	}
	mux := newVehicleTestMux(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestVehicleHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockVehicleSvc{
		getByIDFn: func(_ context.Context, _ uuid.UUID) (*fleet.Vehicle, error) {
			return nil, fleet.ErrNotFound
		},
	}
	mux := newVehicleTestMux(svc)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/vehicles/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestVehicleHandler_UpdateDesignation_Valid(t *testing.T) {
	svc := &mockVehicleSvc{
		updateDesignationFn: func(_ context.Context, _ uuid.UUID, _ fleet.UpdateDesignationInput) (*fleet.Vehicle, error) {
			return &fleet.Vehicle{Designation: fleet.DesignationSalesOnly}, nil
		},
	}
	mux := newVehicleTestMux(svc)
	body, _ := json.Marshal(map[string]string{"designation": "sales_only"})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/vehicles/"+uuid.New().String()+"/designation", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestVehicleHandler_Delete_204(t *testing.T) {
	svc := &mockVehicleSvc{
		deleteFn: func(_ context.Context, _ uuid.UUID, _ string) error {
			return nil
		},
	}
	mux := newVehicleTestMux(svc)
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/vehicles/"+uuid.New().String(), nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", rr.Code)
	}
}

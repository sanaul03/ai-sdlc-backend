package fleet_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/sanaul03/ai-sdlc-backend/internal/fleet"
)

// --- Car Group mock repository ---

type mockCarGroupRepo struct {
	groups map[uuid.UUID]*fleet.CarGroup
}

func newMockCarGroupRepo() *mockCarGroupRepo {
	return &mockCarGroupRepo{groups: make(map[uuid.UUID]*fleet.CarGroup)}
}

func (m *mockCarGroupRepo) Create(_ context.Context, g *fleet.CarGroup) error {
	m.groups[g.ID] = g
	return nil
}

func (m *mockCarGroupRepo) GetByID(_ context.Context, id uuid.UUID) (*fleet.CarGroup, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, fleet.ErrNotFound
	}
	return g, nil
}

func (m *mockCarGroupRepo) List(_ context.Context) ([]*fleet.CarGroup, error) {
	result := make([]*fleet.CarGroup, 0, len(m.groups))
	for _, g := range m.groups {
		result = append(result, g)
	}
	return result, nil
}

func (m *mockCarGroupRepo) Update(_ context.Context, g *fleet.CarGroup) error {
	if _, ok := m.groups[g.ID]; !ok {
		return fleet.ErrNotFound
	}
	existing := m.groups[g.ID]
	existing.Name = g.Name
	existing.Description = g.Description
	existing.SizeCategory = g.SizeCategory
	existing.UpdatedAt = g.UpdatedAt
	existing.UpdatedBy = g.UpdatedBy
	return nil
}

func (m *mockCarGroupRepo) Delete(_ context.Context, id uuid.UUID, deletedBy string) error {
	if _, ok := m.groups[id]; !ok {
		return fleet.ErrNotFound
	}
	delete(m.groups, id)
	return nil
}

// --- Vehicle mock repository ---

type mockVehicleRepo struct {
	vehicles map[uuid.UUID]*fleet.Vehicle
}

func newMockVehicleRepo() *mockVehicleRepo {
	return &mockVehicleRepo{vehicles: make(map[uuid.UUID]*fleet.Vehicle)}
}

func (m *mockVehicleRepo) Create(_ context.Context, v *fleet.Vehicle) error {
	m.vehicles[v.ID] = v
	return nil
}

func (m *mockVehicleRepo) GetByID(_ context.Context, id uuid.UUID) (*fleet.Vehicle, error) {
	v, ok := m.vehicles[id]
	if !ok {
		return nil, fleet.ErrNotFound
	}
	return v, nil
}

func (m *mockVehicleRepo) List(_ context.Context) ([]*fleet.Vehicle, error) {
	result := make([]*fleet.Vehicle, 0, len(m.vehicles))
	for _, v := range m.vehicles {
		result = append(result, v)
	}
	return result, nil
}

func (m *mockVehicleRepo) Update(_ context.Context, v *fleet.Vehicle) error {
	if _, ok := m.vehicles[v.ID]; !ok {
		return fleet.ErrNotFound
	}
	m.vehicles[v.ID] = v
	return nil
}

func (m *mockVehicleRepo) Delete(_ context.Context, id uuid.UUID, _ string) error {
	if _, ok := m.vehicles[id]; !ok {
		return fleet.ErrNotFound
	}
	delete(m.vehicles, id)
	return nil
}

// --- helpers ---

func newRouter(h *fleet.Handler) chi.Router {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

// --- CarGroup handler tests ---

func TestCreateCarGroup_Success(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := `{"name":"Economy Sedan","created_by":"user1"}`
	req := httptest.NewRequest(http.MethodPost, "/car-groups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp fleet.CarGroup
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Name != "Economy Sedan" {
		t.Errorf("expected name 'Economy Sedan', got '%s'", resp.Name)
	}
	if resp.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
}

func TestCreateCarGroup_MissingName(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := `{"created_by":"user1"}`
	req := httptest.NewRequest(http.MethodPost, "/car-groups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetCarGroup_NotFound(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/car-groups/%s", uuid.New()), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestGetCarGroup_Success(t *testing.T) {
	repo := newMockCarGroupRepo()
	h := fleet.NewHandler(repo, newMockVehicleRepo())
	r := newRouter(h)

	// Pre-seed a group
	id := uuid.New()
	now := time.Now().UTC()
	_ = repo.Create(context.Background(), &fleet.CarGroup{
		ID:        id,
		Name:      "Luxury SUV",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "admin",
		UpdatedBy: "admin",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/car-groups/%s", id), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp fleet.CarGroup
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Name != "Luxury SUV" {
		t.Errorf("expected name 'Luxury SUV', got '%s'", resp.Name)
	}
}

func TestListCarGroups_Empty(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/car-groups", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp []*fleet.CarGroup
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty list, got %d items", len(resp))
	}
}

func TestUpdateCarGroup_Success(t *testing.T) {
	repo := newMockCarGroupRepo()
	h := fleet.NewHandler(repo, newMockVehicleRepo())
	r := newRouter(h)

	id := uuid.New()
	now := time.Now().UTC()
	_ = repo.Create(context.Background(), &fleet.CarGroup{
		ID:        id,
		Name:      "Economy",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "user1",
		UpdatedBy: "user1",
	})

	body := `{"name":"Economy Sedan","updated_by":"user2"}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/car-groups/%s", id), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestDeleteCarGroup_Success(t *testing.T) {
	repo := newMockCarGroupRepo()
	h := fleet.NewHandler(repo, newMockVehicleRepo())
	r := newRouter(h)

	id := uuid.New()
	now := time.Now().UTC()
	_ = repo.Create(context.Background(), &fleet.CarGroup{
		ID:        id,
		Name:      "Economy",
		CreatedAt: now,
		UpdatedAt: now,
		CreatedBy: "user1",
		UpdatedBy: "user1",
	})

	body := `{"deleted_by":"admin"}`
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/car-groups/%s", id), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestDeleteCarGroup_NotFound(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := `{"deleted_by":"admin"}`
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/car-groups/%s", uuid.New()), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

// --- Vehicle handler tests ---

func sampleVehicleBody(carGroupID, branchID uuid.UUID) string {
	return fmt.Sprintf(`{
		"car_group_id": "%s",
		"branch_id": "%s",
		"vin": "1HGBH41JXMN109186",
		"licence_plate": "ABC-1234",
		"brand": "Toyota",
		"model": "Corolla",
		"year": 2022,
		"fuel_type": "petrol",
		"transmission_type": "automatic",
		"current_mileage": 5000,
		"status": "available",
		"designation": "rental_only",
		"acquisition_date": "2022-01-15T00:00:00Z",
		"ownership_type": "owned",
		"created_by": "user1"
	}`, carGroupID, branchID)
}

func TestCreateVehicle_Success(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := sampleVehicleBody(uuid.New(), uuid.New())
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp fleet.Vehicle
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.VIN != "1HGBH41JXMN109186" {
		t.Errorf("expected VIN '1HGBH41JXMN109186', got '%s'", resp.VIN)
	}
	if resp.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
}

func TestCreateVehicle_MissingFields(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := `{"brand":"Toyota","created_by":"user1"}`
	req := httptest.NewRequest(http.MethodPost, "/vehicles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestGetVehicle_NotFound(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/vehicles/%s", uuid.New()), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestGetVehicle_Success(t *testing.T) {
	vehicleRepo := newMockVehicleRepo()
	h := fleet.NewHandler(newMockCarGroupRepo(), vehicleRepo)
	r := newRouter(h)

	id := uuid.New()
	now := time.Now().UTC()
	_ = vehicleRepo.Create(context.Background(), &fleet.Vehicle{
		ID:              id,
		CarGroupID:      uuid.New(),
		BranchID:        uuid.New(),
		VIN:             "TEST123",
		LicencePlate:    "XYZ-999",
		Brand:           "BMW",
		Model:           "X5",
		Year:            2023,
		FuelType:        fleet.FuelTypePetrol,
		TransmissionType: fleet.TransmissionTypeAutomatic,
		Status:          fleet.VehicleStatusAvailable,
		Designation:     fleet.VehicleDesignationRentalOnly,
		AcquisitionDate: now,
		OwnershipType:   fleet.OwnershipTypeOwned,
		CreatedAt:       now,
		UpdatedAt:       now,
		CreatedBy:       "admin",
		UpdatedBy:       "admin",
	})

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/vehicles/%s", id), nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}

	var resp fleet.Vehicle
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.VIN != "TEST123" {
		t.Errorf("expected VIN 'TEST123', got '%s'", resp.VIN)
	}
}

func TestListVehicles_Empty(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/vehicles", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp []*fleet.Vehicle
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp) != 0 {
		t.Errorf("expected empty list, got %d items", len(resp))
	}
}

func TestDeleteVehicle_NotFound(t *testing.T) {
	h := fleet.NewHandler(newMockCarGroupRepo(), newMockVehicleRepo())
	r := newRouter(h)

	body := `{"deleted_by":"admin"}`
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/vehicles/%s", uuid.New()), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

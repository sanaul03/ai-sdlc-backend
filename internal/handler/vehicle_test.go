package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanaul03/ai-sdlc-backend/internal/handler"
	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// ---- in-memory Vehicle repository stub ----

type stubVehicleRepo struct {
	data   map[int64]*model.Vehicle
	nextID int64
}

func newStubVehicleRepo() *stubVehicleRepo {
	return &stubVehicleRepo{data: make(map[int64]*model.Vehicle), nextID: 1}
}

func (s *stubVehicleRepo) Create(_ context.Context, req model.CreateVehicleRequest) (*model.Vehicle, error) {
	v := &model.Vehicle{
		ID: s.nextID, CompanyID: req.CompanyID, DepotID: req.DepotID,
		VehicleTypeID: req.VehicleTypeID, RegistrationNumber: req.RegistrationNumber,
		ChassisNumber: req.ChassisNumber, EngineNumber: req.EngineNumber,
		ManufactureYear: req.ManufactureYear, Color: req.Color,
		Status: model.VehicleStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.data[s.nextID] = v
	s.nextID++
	return v, nil
}

func (s *stubVehicleRepo) GetByID(_ context.Context, id int64) (*model.Vehicle, error) {
	v, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	return v, nil
}

func (s *stubVehicleRepo) ListByCompany(_ context.Context, companyID int64) ([]model.Vehicle, error) {
	var list []model.Vehicle
	for _, v := range s.data {
		if v.CompanyID == companyID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (s *stubVehicleRepo) Update(_ context.Context, id int64, req model.UpdateVehicleRequest) (*model.Vehicle, error) {
	v, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	if req.Status != nil {
		v.Status = *req.Status
	}
	v.UpdatedAt = time.Now()
	return v, nil
}

func (s *stubVehicleRepo) Delete(_ context.Context, id int64) error {
	if _, ok := s.data[id]; !ok {
		return errNotFound
	}
	delete(s.data, id)
	return nil
}

// ---- Vehicle handler tests ----

func TestVehicleHandler_Create(t *testing.T) {
	h := handler.NewVehicleHandler(newStubVehicleRepo())

	body := model.CreateVehicleRequest{
		CompanyID: 1, VehicleTypeID: 2, RegistrationNumber: "ABC-1234",
	}
	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/vehicles", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}
	var got model.Vehicle
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.RegistrationNumber != body.RegistrationNumber {
		t.Errorf("want reg# %q, got %q", body.RegistrationNumber, got.RegistrationNumber)
	}
	if got.Status != model.VehicleStatusActive {
		t.Errorf("want status %q, got %q", model.VehicleStatusActive, got.Status)
	}
}

func TestVehicleHandler_Create_MissingFields(t *testing.T) {
	h := handler.NewVehicleHandler(newStubVehicleRepo())

	w := httptest.NewRecorder()
	// missing registration_number
	h.Create(w, newReq(t, http.MethodPost, "/vehicles",
		model.CreateVehicleRequest{CompanyID: 1, VehicleTypeID: 2}))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestVehicleHandler_GetByID_NotFound(t *testing.T) {
	h := handler.NewVehicleHandler(newStubVehicleRepo())

	r := httptest.NewRequest(http.MethodGet, "/vehicles/999", nil)
	r.SetPathValue("id", "999")
	w := httptest.NewRecorder()
	h.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestVehicleHandler_ListByCompany(t *testing.T) {
	repo := newStubVehicleRepo()
	h := handler.NewVehicleHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateVehicleRequest{
		CompanyID: 3, VehicleTypeID: 1, RegistrationNumber: "V1",
	})
	_, _ = repo.Create(context.Background(), model.CreateVehicleRequest{
		CompanyID: 3, VehicleTypeID: 1, RegistrationNumber: "V2",
	})
	_, _ = repo.Create(context.Background(), model.CreateVehicleRequest{
		CompanyID: 9, VehicleTypeID: 1, RegistrationNumber: "V3",
	})

	r := httptest.NewRequest(http.MethodGet, "/companies/3/vehicles", nil)
	r.SetPathValue("id", "3")
	w := httptest.NewRecorder()
	h.ListByCompany(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []model.Vehicle
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("want 2 vehicles for company 3, got %d", len(got))
	}
}

func TestVehicleHandler_Update_Status(t *testing.T) {
	repo := newStubVehicleRepo()
	h := handler.NewVehicleHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateVehicleRequest{
		CompanyID: 1, VehicleTypeID: 1, RegistrationNumber: "XY-99",
	})

	newStatus := model.VehicleStatusMaintenance
	r := httptest.NewRequest(http.MethodPut, "/vehicles/1",
		newReqBodyStatus(t, newStatus))
	r.Header.Set("Content-Type", "application/json")
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Update(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got model.Vehicle
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Status != newStatus {
		t.Errorf("want status %q, got %q", newStatus, got.Status)
	}
}

func newReqBodyStatus(t *testing.T, status string) *bytes.Buffer {
	t.Helper()
	b, _ := json.Marshal(model.UpdateVehicleRequest{Status: &status})
	return bytes.NewBuffer(b)
}
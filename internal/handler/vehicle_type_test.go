package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanaul03/ai-sdlc-backend/internal/handler"
	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

// ---- in-memory VehicleType repository stub ----

type stubVehicleTypeRepo struct {
	data   map[int64]*model.VehicleType
	nextID int64
}

func newStubVehicleTypeRepo() *stubVehicleTypeRepo {
	return &stubVehicleTypeRepo{data: make(map[int64]*model.VehicleType), nextID: 1}
}

func (s *stubVehicleTypeRepo) Create(_ context.Context, req model.CreateVehicleTypeRequest) (*model.VehicleType, error) {
	vt := &model.VehicleType{
		ID: s.nextID, CategoryID: req.CategoryID, Name: req.Name, Code: req.Code,
		Capacity: req.Capacity, Description: req.Description,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.data[s.nextID] = vt
	s.nextID++
	return vt, nil
}

func (s *stubVehicleTypeRepo) GetByID(_ context.Context, id int64) (*model.VehicleType, error) {
	vt, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	return vt, nil
}

func (s *stubVehicleTypeRepo) ListByCategory(_ context.Context, categoryID int64) ([]model.VehicleType, error) {
	var list []model.VehicleType
	for _, vt := range s.data {
		if vt.CategoryID == categoryID {
			list = append(list, *vt)
		}
	}
	return list, nil
}

func (s *stubVehicleTypeRepo) Update(_ context.Context, id int64, req model.UpdateVehicleTypeRequest) (*model.VehicleType, error) {
	vt, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	if req.Name != nil {
		vt.Name = *req.Name
	}
	vt.UpdatedAt = time.Now()
	return vt, nil
}

func (s *stubVehicleTypeRepo) Delete(_ context.Context, id int64) error {
	if _, ok := s.data[id]; !ok {
		return errNotFound
	}
	delete(s.data, id)
	return nil
}

// ---- VehicleType handler tests ----

func TestVehicleTypeHandler_Create(t *testing.T) {
	h := handler.NewVehicleTypeHandler(newStubVehicleTypeRepo())

	cap := 50
	body := model.CreateVehicleTypeRequest{CategoryID: 1, Name: "Minibus", Code: "MB", Capacity: &cap}
	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/vehicle-types", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}
	var got model.VehicleType
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Name != body.Name {
		t.Errorf("want name %q, got %q", body.Name, got.Name)
	}
}

func TestVehicleTypeHandler_Create_MissingFields(t *testing.T) {
	h := handler.NewVehicleTypeHandler(newStubVehicleTypeRepo())

	w := httptest.NewRecorder()
	// missing category_id
	h.Create(w, newReq(t, http.MethodPost, "/vehicle-types",
		model.CreateVehicleTypeRequest{Name: "Bus", Code: "BUS"}))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestVehicleTypeHandler_ListByCategory(t *testing.T) {
	repo := newStubVehicleTypeRepo()
	h := handler.NewVehicleTypeHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateVehicleTypeRequest{CategoryID: 2, Name: "T1", Code: "T1"})
	_, _ = repo.Create(context.Background(), model.CreateVehicleTypeRequest{CategoryID: 2, Name: "T2", Code: "T2"})
	_, _ = repo.Create(context.Background(), model.CreateVehicleTypeRequest{CategoryID: 7, Name: "Other", Code: "OT"})

	r := httptest.NewRequest(http.MethodGet, "/vehicle-categories/2/vehicle-types", nil)
	r.SetPathValue("id", "2")
	w := httptest.NewRecorder()
	h.ListByCategory(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []model.VehicleType
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("want 2 types for category 2, got %d", len(got))
	}
}

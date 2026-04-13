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

// ---- in-memory VehicleCategory repository stub ----

type stubVehicleCategoryRepo struct {
	data   map[int64]*model.VehicleCategory
	nextID int64
}

func newStubVehicleCategoryRepo() *stubVehicleCategoryRepo {
	return &stubVehicleCategoryRepo{data: make(map[int64]*model.VehicleCategory), nextID: 1}
}

func (s *stubVehicleCategoryRepo) Create(_ context.Context, req model.CreateVehicleCategoryRequest) (*model.VehicleCategory, error) {
	vc := &model.VehicleCategory{
		ID: s.nextID, Name: req.Name, Code: req.Code, Description: req.Description,
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.data[s.nextID] = vc
	s.nextID++
	return vc, nil
}

func (s *stubVehicleCategoryRepo) GetByID(_ context.Context, id int64) (*model.VehicleCategory, error) {
	vc, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	return vc, nil
}

func (s *stubVehicleCategoryRepo) List(_ context.Context) ([]model.VehicleCategory, error) {
	var list []model.VehicleCategory
	for _, vc := range s.data {
		list = append(list, *vc)
	}
	return list, nil
}

func (s *stubVehicleCategoryRepo) Update(_ context.Context, id int64, req model.UpdateVehicleCategoryRequest) (*model.VehicleCategory, error) {
	vc, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	if req.Name != nil {
		vc.Name = *req.Name
	}
	vc.UpdatedAt = time.Now()
	return vc, nil
}

func (s *stubVehicleCategoryRepo) Delete(_ context.Context, id int64) error {
	if _, ok := s.data[id]; !ok {
		return errNotFound
	}
	delete(s.data, id)
	return nil
}

// ---- VehicleCategory handler tests ----

func TestVehicleCategoryHandler_Create(t *testing.T) {
	h := handler.NewVehicleCategoryHandler(newStubVehicleCategoryRepo())

	body := model.CreateVehicleCategoryRequest{Name: "Bus", Code: "BUS"}
	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/vehicle-categories", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}
	var got model.VehicleCategory
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Code != "BUS" {
		t.Errorf("want code BUS, got %s", got.Code)
	}
}

func TestVehicleCategoryHandler_List(t *testing.T) {
	repo := newStubVehicleCategoryRepo()
	h := handler.NewVehicleCategoryHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateVehicleCategoryRequest{Name: "Bus", Code: "BUS"})
	_, _ = repo.Create(context.Background(), model.CreateVehicleCategoryRequest{Name: "Truck", Code: "TRK"})

	w := httptest.NewRecorder()
	h.List(w, httptest.NewRequest(http.MethodGet, "/vehicle-categories", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []model.VehicleCategory
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("want 2 categories, got %d", len(got))
	}
}

func TestVehicleCategoryHandler_Delete(t *testing.T) {
	repo := newStubVehicleCategoryRepo()
	h := handler.NewVehicleCategoryHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateVehicleCategoryRequest{Name: "Car", Code: "CAR"})

	r := httptest.NewRequest(http.MethodDelete, "/vehicle-categories/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Delete(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", w.Code)
	}
}

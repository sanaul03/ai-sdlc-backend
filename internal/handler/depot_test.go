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

// ---- in-memory Depot repository stub ----

type stubDepotRepo struct {
	data   map[int64]*model.Depot
	nextID int64
}

func newStubDepotRepo() *stubDepotRepo {
	return &stubDepotRepo{data: make(map[int64]*model.Depot), nextID: 1}
}

func (s *stubDepotRepo) Create(_ context.Context, req model.CreateDepotRequest) (*model.Depot, error) {
	d := &model.Depot{
		ID: s.nextID, CompanyID: req.CompanyID, Name: req.Name, Code: req.Code,
		Address: req.Address, Latitude: req.Latitude, Longitude: req.Longitude,
		Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.data[s.nextID] = d
	s.nextID++
	return d, nil
}

func (s *stubDepotRepo) GetByID(_ context.Context, id int64) (*model.Depot, error) {
	d, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	return d, nil
}

func (s *stubDepotRepo) ListByCompany(_ context.Context, companyID int64) ([]model.Depot, error) {
	var list []model.Depot
	for _, d := range s.data {
		if d.CompanyID == companyID {
			list = append(list, *d)
		}
	}
	return list, nil
}

func (s *stubDepotRepo) Update(_ context.Context, id int64, req model.UpdateDepotRequest) (*model.Depot, error) {
	d, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	if req.Name != nil {
		d.Name = *req.Name
	}
	d.UpdatedAt = time.Now()
	return d, nil
}

func (s *stubDepotRepo) Delete(_ context.Context, id int64) error {
	if _, ok := s.data[id]; !ok {
		return errNotFound
	}
	delete(s.data, id)
	return nil
}

// ---- Depot handler tests ----

func TestDepotHandler_Create(t *testing.T) {
	h := handler.NewDepotHandler(newStubDepotRepo())

	body := model.CreateDepotRequest{CompanyID: 1, Name: "Main Depot", Code: "MD"}
	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/depots", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}
	var got model.Depot
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Name != body.Name {
		t.Errorf("want name %q, got %q", body.Name, got.Name)
	}
}

func TestDepotHandler_Create_MissingFields(t *testing.T) {
	h := handler.NewDepotHandler(newStubDepotRepo())

	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/depots", model.CreateDepotRequest{Name: "No Company"}))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestDepotHandler_GetByID_NotFound(t *testing.T) {
	h := handler.NewDepotHandler(newStubDepotRepo())

	r := httptest.NewRequest(http.MethodGet, "/depots/99", nil)
	r.SetPathValue("id", "99")
	w := httptest.NewRecorder()
	h.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestDepotHandler_ListByCompany(t *testing.T) {
	repo := newStubDepotRepo()
	h := handler.NewDepotHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateDepotRequest{CompanyID: 5, Name: "D1", Code: "D1"})
	_, _ = repo.Create(context.Background(), model.CreateDepotRequest{CompanyID: 5, Name: "D2", Code: "D2"})
	_, _ = repo.Create(context.Background(), model.CreateDepotRequest{CompanyID: 9, Name: "Other", Code: "OT"})

	r := httptest.NewRequest(http.MethodGet, "/companies/5/depots", nil)
	r.SetPathValue("id", "5")
	w := httptest.NewRecorder()
	h.ListByCompany(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []model.Depot
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("want 2 depots for company 5, got %d", len(got))
	}
}

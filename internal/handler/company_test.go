package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sanaul03/ai-sdlc-backend/internal/handler"
	"github.com/sanaul03/ai-sdlc-backend/internal/model"
)

var errNotFound = errors.New("not found")

// ---- in-memory Company repository stub ----

type stubCompanyRepo struct {
	data   map[int64]*model.Company
	nextID int64
}

func newStubCompanyRepo() *stubCompanyRepo {
	return &stubCompanyRepo{data: make(map[int64]*model.Company), nextID: 1}
}

func (s *stubCompanyRepo) Create(_ context.Context, req model.CreateCompanyRequest) (*model.Company, error) {
	c := &model.Company{
		ID: s.nextID, Name: req.Name, Code: req.Code,
		Address: req.Address, Phone: req.Phone, Email: req.Email,
		Status: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	s.data[s.nextID] = c
	s.nextID++
	return c, nil
}

func (s *stubCompanyRepo) GetByID(_ context.Context, id int64) (*model.Company, error) {
	c, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	return c, nil
}

func (s *stubCompanyRepo) List(_ context.Context) ([]model.Company, error) {
	var list []model.Company
	for _, c := range s.data {
		list = append(list, *c)
	}
	return list, nil
}

func (s *stubCompanyRepo) Update(_ context.Context, id int64, req model.UpdateCompanyRequest) (*model.Company, error) {
	c, ok := s.data[id]
	if !ok {
		return nil, errNotFound
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	c.UpdatedAt = time.Now()
	return c, nil
}

func (s *stubCompanyRepo) Delete(_ context.Context, id int64) error {
	if _, ok := s.data[id]; !ok {
		return errNotFound
	}
	delete(s.data, id)
	return nil
}

// ---- helper: build a request with optional JSON body ----

func newReq(t *testing.T, method, target string, body any) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			t.Fatalf("encode body: %v", err)
		}
	}
	r := httptest.NewRequest(method, target, &buf)
	r.Header.Set("Content-Type", "application/json")
	return r
}

// ---- Company handler tests ----

func TestCompanyHandler_Create(t *testing.T) {
	h := handler.NewCompanyHandler(newStubCompanyRepo())

	body := model.CreateCompanyRequest{Name: "Acme Corp", Code: "ACME"}
	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/companies", body))

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d", w.Code)
	}

	var got model.Company
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Name != body.Name || got.Code != body.Code {
		t.Errorf("want %+v, got %+v", body, got)
	}
}

func TestCompanyHandler_Create_MissingFields(t *testing.T) {
	h := handler.NewCompanyHandler(newStubCompanyRepo())

	w := httptest.NewRecorder()
	h.Create(w, newReq(t, http.MethodPost, "/companies", model.CreateCompanyRequest{Name: "Only Name"}))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCompanyHandler_GetByID_NotFound(t *testing.T) {
	h := handler.NewCompanyHandler(newStubCompanyRepo())

	r := httptest.NewRequest(http.MethodGet, "/companies/99", nil)
	r.SetPathValue("id", "99")
	w := httptest.NewRecorder()
	h.GetByID(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestCompanyHandler_GetByID_Found(t *testing.T) {
	repo := newStubCompanyRepo()
	h := handler.NewCompanyHandler(repo)

	// pre-populate
	created, _ := repo.Create(context.Background(), model.CreateCompanyRequest{Name: "Test", Code: "TST"})

	r := httptest.NewRequest(http.MethodGet, "/companies/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.GetByID(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got model.Company
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.ID != created.ID {
		t.Errorf("want ID %d, got %d", created.ID, got.ID)
	}
}

func TestCompanyHandler_List(t *testing.T) {
	repo := newStubCompanyRepo()
	h := handler.NewCompanyHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateCompanyRequest{Name: "Alpha", Code: "A"})
	_, _ = repo.Create(context.Background(), model.CreateCompanyRequest{Name: "Beta", Code: "B"})

	w := httptest.NewRecorder()
	h.List(w, httptest.NewRequest(http.MethodGet, "/companies", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got []model.Company
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 {
		t.Errorf("want 2 companies, got %d", len(got))
	}
}

func TestCompanyHandler_Update(t *testing.T) {
	repo := newStubCompanyRepo()
	h := handler.NewCompanyHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateCompanyRequest{Name: "Old Name", Code: "OC"})

	newName := "New Name"
	r := httptest.NewRequest(http.MethodPut, "/companies/1",
		bytes.NewBufferString(`{"name":"New Name"}`))
	r.Header.Set("Content-Type", "application/json")
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Update(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var got model.Company
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.Name != newName {
		t.Errorf("want name %q, got %q", newName, got.Name)
	}
}

func TestCompanyHandler_Delete(t *testing.T) {
	repo := newStubCompanyRepo()
	h := handler.NewCompanyHandler(repo)

	_, _ = repo.Create(context.Background(), model.CreateCompanyRequest{Name: "ToDelete", Code: "TD"})

	r := httptest.NewRequest(http.MethodDelete, "/companies/1", nil)
	r.SetPathValue("id", "1")
	w := httptest.NewRecorder()
	h.Delete(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d", w.Code)
	}
}

func TestCompanyHandler_InvalidID(t *testing.T) {
	h := handler.NewCompanyHandler(newStubCompanyRepo())

	r := httptest.NewRequest(http.MethodGet, "/companies/abc", nil)
	r.SetPathValue("id", "abc")
	w := httptest.NewRecorder()
	h.GetByID(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

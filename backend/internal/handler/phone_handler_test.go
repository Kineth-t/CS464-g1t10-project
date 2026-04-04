package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// newTestPhoneHandler wires a PhoneHandler backed by in-memory mocks.
// Seed phones are pre-loaded into the repo before returning.
func newTestPhoneHandler(seed ...model.Phone) *PhoneHandler {
	repo := newMockPhoneRepo()
	for _, p := range seed {
		repo.Create(p)
	}
	cache := repository.NewPhoneCache(nil) // no-op cache
	svc := service.NewPhoneService(repo, cache)
	return NewPhoneHandler(svc)
}

func TestListPhones_Empty(t *testing.T) {
	h := newTestPhoneHandler()
	req := httptest.NewRequest(http.MethodGet, "/phones", nil)
	w := httptest.NewRecorder()

	h.ListPhones(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var phones []model.Phone
	if err := json.NewDecoder(w.Body).Decode(&phones); err != nil {
		t.Fatalf("response is not valid JSON array: %v", err)
	}
	if len(phones) != 0 {
		t.Errorf("expected empty list, got %d items", len(phones))
	}
}

func TestListPhones_ReturnsAll(t *testing.T) {
	h := newTestPhoneHandler(
		model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 999.0, Stock: 10},
		model.Phone{Brand: "Samsung", Model: "Galaxy S24", Price: 899.0, Stock: 5},
	)
	req := httptest.NewRequest(http.MethodGet, "/phones", nil)
	w := httptest.NewRecorder()

	h.ListPhones(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var phones []model.Phone
	if err := json.NewDecoder(w.Body).Decode(&phones); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if len(phones) != 2 {
		t.Errorf("expected 2 phones, got %d", len(phones))
	}
}

func TestGetPhone_Found(t *testing.T) {
	h := newTestPhoneHandler(
		model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 999.0, Stock: 10},
	)
	req := httptest.NewRequest(http.MethodGet, "/phones/1", nil)
	w := httptest.NewRecorder()

	h.GetPhone(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var phone model.Phone
	if err := json.NewDecoder(w.Body).Decode(&phone); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if phone.Brand != "Apple" {
		t.Errorf("expected Brand=Apple, got %q", phone.Brand)
	}
}

func TestGetPhone_NotFound(t *testing.T) {
	h := newTestPhoneHandler()
	req := httptest.NewRequest(http.MethodGet, "/phones/99", nil)
	w := httptest.NewRecorder()

	h.GetPhone(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetPhone_InvalidID(t *testing.T) {
	h := newTestPhoneHandler()
	req := httptest.NewRequest(http.MethodGet, "/phones/abc", nil)
	w := httptest.NewRecorder()

	h.GetPhone(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreatePhone_Success(t *testing.T) {
	h := newTestPhoneHandler()
	body, _ := json.Marshal(model.Phone{Brand: "Google", Model: "Pixel 9", Price: 799.0, Stock: 20})
	req := httptest.NewRequest(http.MethodPost, "/phones", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreatePhone(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var created model.Phone
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if created.ID == 0 {
		t.Error("expected a non-zero ID on the created phone")
	}
}

func TestCreatePhone_ZeroPrice(t *testing.T) {
	h := newTestPhoneHandler()
	body, _ := json.Marshal(model.Phone{Brand: "Google", Model: "Pixel 9", Price: 0, Stock: 20})
	req := httptest.NewRequest(http.MethodPost, "/phones", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.CreatePhone(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero price, got %d", w.Code)
	}
}

func TestCreatePhone_BadJSON(t *testing.T) {
	h := newTestPhoneHandler()
	req := httptest.NewRequest(http.MethodPost, "/phones", bytes.NewReader([]byte("{invalid")))
	w := httptest.NewRecorder()

	h.CreatePhone(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed JSON, got %d", w.Code)
	}
}

func TestDeletePhone_Success(t *testing.T) {
	h := newTestPhoneHandler(
		model.Phone{Brand: "OnePlus", Model: "12", Price: 699.0, Stock: 8},
	)
	req := httptest.NewRequest(http.MethodDelete, "/phones/1", nil)
	w := httptest.NewRecorder()

	h.DeletePhone(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestDeletePhone_NotFound(t *testing.T) {
	h := newTestPhoneHandler()
	req := httptest.NewRequest(http.MethodDelete, "/phones/999", nil)
	w := httptest.NewRecorder()

	h.DeletePhone(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

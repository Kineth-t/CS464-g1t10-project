package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// newTestCartHandler wires a CartHandler backed by in-memory mocks.
// Seed phones are pre-loaded into the phone repo before returning.
func newTestCartHandler(seed ...model.Phone) *CartHandler {
	phoneRepo := newMockPhoneRepo()
	for _, p := range seed {
		phoneRepo.Create(p)
	}
	svc := service.NewCartService(newMockCartRepo(), phoneRepo)
	return NewCartHandler(svc)
}

// withUserID injects a user ID into the request context, bypassing auth middleware.
func withUserID(r *http.Request, userID int) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserIDKey, userID)
	return r.WithContext(ctx)
}

func TestGetCart_NoCart(t *testing.T) {
	h := newTestCartHandler()
	req := withUserID(httptest.NewRequest(http.MethodGet, "/cart", nil), 1)
	w := httptest.NewRecorder()

	h.GetCart(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 when no cart exists, got %d", w.Code)
	}
}

func TestGetCart_WithItems(t *testing.T) {
	h := newTestCartHandler(
		model.Phone{Brand: "Apple", Model: "iPhone 15", Price: 999.0, Stock: 10},
	)

	// Add an item first
	addBody, _ := json.Marshal(map[string]int{"phone_id": 1, "quantity": 1})
	h.AddToCart(httptest.NewRecorder(), withUserID(
		httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(addBody)), 1,
	))

	// Fetch the cart
	req := withUserID(httptest.NewRequest(http.MethodGet, "/cart", nil), 1)
	w := httptest.NewRecorder()
	h.GetCart(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var cart model.Cart
	if err := json.NewDecoder(w.Body).Decode(&cart); err != nil {
		t.Fatalf("could not decode cart response: %v", err)
	}
	if len(cart.Items) == 0 {
		t.Error("expected at least one item in cart")
	}
}

func TestAddToCart_Success(t *testing.T) {
	h := newTestCartHandler(
		model.Phone{Brand: "Samsung", Model: "Galaxy S24", Price: 899.0, Stock: 5},
	)
	body, _ := json.Marshal(map[string]int{"phone_id": 1, "quantity": 2})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body)), 42)
	w := httptest.NewRecorder()

	h.AddToCart(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var item model.CartItem
	if err := json.NewDecoder(w.Body).Decode(&item); err != nil {
		t.Fatalf("could not decode cart item: %v", err)
	}
	if item.PhoneID != 1 {
		t.Errorf("expected phone_id=1, got %d", item.PhoneID)
	}
	if item.Quantity != 2 {
		t.Errorf("expected quantity=2, got %d", item.Quantity)
	}
}

func TestAddToCart_PhoneNotFound(t *testing.T) {
	h := newTestCartHandler() // no phones seeded
	body, _ := json.Marshal(map[string]int{"phone_id": 999, "quantity": 1})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body)), 1)
	w := httptest.NewRecorder()

	h.AddToCart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-existent phone, got %d", w.Code)
	}
}

func TestAddToCart_InsufficientStock(t *testing.T) {
	h := newTestCartHandler(
		model.Phone{Brand: "Nokia", Model: "3310", Price: 49.0, Stock: 1},
	)
	body, _ := json.Marshal(map[string]int{"phone_id": 1, "quantity": 5})
	req := withUserID(httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(body)), 1)
	w := httptest.NewRecorder()

	h.AddToCart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for insufficient stock, got %d", w.Code)
	}
}

func TestAddToCart_BadJSON(t *testing.T) {
	h := newTestCartHandler()
	req := withUserID(httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader([]byte("{bad}"))), 1)
	w := httptest.NewRecorder()

	h.AddToCart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed JSON, got %d", w.Code)
	}
}

func TestRemoveFromCart_Success(t *testing.T) {
	h := newTestCartHandler(
		model.Phone{Brand: "Sony", Model: "Xperia 5", Price: 599.0, Stock: 10},
	)

	// Add an item and capture its ID
	addBody, _ := json.Marshal(map[string]int{"phone_id": 1, "quantity": 1})
	addW := httptest.NewRecorder()
	h.AddToCart(addW, withUserID(
		httptest.NewRequest(http.MethodPost, "/cart", bytes.NewReader(addBody)), 1,
	))

	var item model.CartItem
	json.NewDecoder(addW.Body).Decode(&item)

	// Remove the item
	req := withUserID(
		httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/cart/%d", item.ID), nil), 1,
	)
	w := httptest.NewRecorder()
	h.RemoveFromCart(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRemoveFromCart_InvalidID(t *testing.T) {
	h := newTestCartHandler()
	req := withUserID(httptest.NewRequest(http.MethodDelete, "/cart/abc", nil), 1)
	w := httptest.NewRecorder()

	h.RemoveFromCart(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-numeric ID, got %d", w.Code)
	}
}

func TestRemoveFromCart_NotFound(t *testing.T) {
	h := newTestCartHandler()
	req := withUserID(httptest.NewRequest(http.MethodDelete, "/cart/999", nil), 1)
	w := httptest.NewRecorder()

	h.RemoveFromCart(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404 for missing item, got %d", w.Code)
	}
}

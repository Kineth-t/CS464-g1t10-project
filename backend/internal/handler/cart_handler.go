package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	cart, err := h.service.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	var body struct {
		PhoneID  int `json:"phone_id"`
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	item, err := h.service.AddToCart(userID, body.PhoneID, body.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	idStr := strings.TrimPrefix(r.URL.Path, "/cart/")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}
	if err := h.service.RemoveFromCart(userID, itemID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	if err := h.service.Checkout(userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "checkout successful"})
}
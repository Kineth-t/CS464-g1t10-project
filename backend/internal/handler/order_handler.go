package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// OrderHandler handles HTTP requests for order retrieval
type OrderHandler struct {
	service *service.OrderService
}

// Constructor
func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// GetOrders returns all orders for the logged in user
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	orders, err := h.service.GetOrders(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrder returns a single order — verifies it belongs to the logged in user
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Order ID is a Stripe pi_xxx string, not an integer
	orderID := strings.TrimPrefix(r.URL.Path, "/orders/")
	if orderID == "" {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(userID, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
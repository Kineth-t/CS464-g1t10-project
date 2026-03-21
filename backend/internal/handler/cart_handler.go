package handler

import (
	"encoding/json"  // For JSON encoding/decoding
	"net/http"       // HTTP server utilities
	"strconv"        // Convert string → int
	"strings"        // String manipulation

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// CartHandler handles all cart-related HTTP requests
type CartHandler struct {
	service *service.CartService // Business logic layer for cart
}

// Constructor to create a new CartHandler
func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

// GetCart retrieves the current user's cart
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request context (set by authentication middleware)
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Call service layer to fetch cart
	cart, err := h.service.GetCart(userID)
	if err != nil {
		// If cart not found or error occurs
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return cart as JSON
	json.NewEncoder(w).Encode(cart)
}

// AddToCart adds an item to the user's cart
func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Temporary struct to store request body
	var body struct {
		PhoneID  int `json:"phone_id"`  // ID of the phone/product
		Quantity int `json:"quantity"`  // Number of items
	}

	// Decode JSON request into struct
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Call service to add item into cart
	item, err := h.service.AddToCart(userID, body.PhoneID, body.Quantity)
	if err != nil {
		// Handle validation errors (e.g., invalid phone ID, stock issues)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return 201 Created since a new item is added
	w.WriteHeader(http.StatusCreated)

	// Return the added cart item
	json.NewEncoder(w).Encode(item)
}

// RemoveFromCart removes a specific item from the cart
func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Extract item ID from URL path (e.g., "/cart/3" → "3")
	idStr := strings.TrimPrefix(r.URL.Path, "/cart/")

	// Convert string ID to integer
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	// Call service to remove item
	if err := h.service.RemoveFromCart(userID, itemID); err != nil {
		// If item not found or doesn't belong to user
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return 204 No Content (successful deletion, no response body)
	w.WriteHeader(http.StatusNoContent)
}

// Checkout processes the user's cart into an order
func (h *CartHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Call service to perform checkout logic
	if err := h.service.Checkout(userID); err != nil {
		// If checkout fails (e.g., empty cart, insufficient stock)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success message
	json.NewEncoder(w).Encode(map[string]string{
		"message": "checkout successful",
	})
}
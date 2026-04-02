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
//
// @Summary      Get the current user's cart
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  model.Cart
// @Failure      401  {string}  string "unauthorized"
// @Failure      404  {string}  string "cart not found"
// @Router       /cart [get]
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
//
// @Summary      Add a phone to the cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body      object{phone_id=int,quantity=int} true "Item to add"
// @Success      201  {object}  model.CartItem
// @Failure      400  {string}  string "invalid body or phone not found"
// @Failure      401  {string}  string "unauthorized"
// @Router       /cart [post]
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
//
// @Summary      Remove an item from the cart
// @Tags         cart
// @Security     BearerAuth
// @Param        id  path      int  true  "Cart item ID"
// @Success      204 {string}  string ""
// @Failure      400 {string}  string "invalid item id"
// @Failure      401 {string}  string "unauthorized"
// @Failure      404 {string}  string "item not found"
// @Router       /cart/{id} [delete]
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

// NOTE: The /cart/checkout route has been intentionally removed.
// Stock deduction and cart checkout are handled atomically inside
// PaymentService.ProcessPayment after a successful Stripe charge,
// preventing double-deduction if both endpoints were called.
package handler

import (
	"encoding/json"  // For JSON encoding/decoding
	"net/http"       // HTTP server utilities
	"strconv"        // Convert string → int
	"strings"        // String manipulation

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// PhoneHandler handles all phone/product-related HTTP requests
type PhoneHandler struct {
	service *service.PhoneService // Business logic layer
}

// Constructor to create a new PhoneHandler
func NewPhoneHandler(s *service.PhoneService) *PhoneHandler {
	return &PhoneHandler{service: s}
}

// ListPhones returns all phones
func (h *PhoneHandler) ListPhones(w http.ResponseWriter, r *http.Request) {
	// Directly call service and return result as JSON
	json.NewEncoder(w).Encode(h.service.ListPhones())
}

// GetPhone returns a single phone by ID
func (h *PhoneHandler) GetPhone(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL (e.g., "/phones/3" → "3")
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Fetch phone from service
	phone, err := h.service.GetPhone(id)
	if err != nil {
		// If phone not found
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return phone as JSON
	json.NewEncoder(w).Encode(phone)
}

// CreatePhone adds a new phone
func (h *PhoneHandler) CreatePhone(w http.ResponseWriter, r *http.Request) {
	var p model.Phone

	// Decode JSON request body into Phone struct
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Call service to create phone
	created, err := h.service.CreatePhone(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return 201 Created
	w.WriteHeader(http.StatusCreated)

	// Return created phone
	json.NewEncoder(w).Encode(created)
}

// UpdatePhone updates an existing phone
func (h *PhoneHandler) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var p model.Phone

	// Decode updated data from request body
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Set the ID from URL into the struct
	p.ID = id

	// Call service to update
	if err := h.service.UpdatePhone(p); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return updated phone
	json.NewEncoder(w).Encode(p)
}

// DeletePhone removes a phone by ID
func (h *PhoneHandler) DeletePhone(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/phones/"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// Call service to delete
	if err := h.service.DeletePhone(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return 204 No Content (successful deletion)
	w.WriteHeader(http.StatusNoContent)
}

// PurchasePhone handles direct purchase (without cart)
func (h *PhoneHandler) PurchasePhone(w http.ResponseWriter, r *http.Request) {
	// Temporary struct for request body
	var body struct {
		PhoneID  int `json:"phone_id"`  // ID of phone to purchase
		Quantity int `json:"quantity"`  // Quantity to buy
	}

	// Decode JSON request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Call service to process purchase
	if err := h.service.PurchasePhone(body.PhoneID, body.Quantity); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success message
	json.NewEncoder(w).Encode(map[string]string{
		"message": "purchase successful",
	})
}
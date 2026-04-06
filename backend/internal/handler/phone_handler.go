package handler

import (
	"encoding/json"  // For JSON encoding/decoding
	"net/http"       // HTTP server utilities
	"strconv"        // Convert string → int
	"strings"        // String manipulation

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// PhoneHandler handles all phone/product-related HTTP requests
type PhoneHandler struct {
	service *service.PhoneService // Business logic layer
	audit   *service.AuditService
}

// Constructor to create a new PhoneHandler
func NewPhoneHandler(s *service.PhoneService) *PhoneHandler {
	return &PhoneHandler{service: s}
}

// SetAudit attaches an audit service for event logging.
func (h *PhoneHandler) SetAudit(s *service.AuditService) { h.audit = s }

// ListPhones returns all phones
//
// @Summary      List all phones
// @Tags         phones
// @Produce      json
// @Success      200  {array}   model.Phone
// @Router       /phones [get]
func (h *PhoneHandler) ListPhones(w http.ResponseWriter, r *http.Request) {
    phones, err := h.service.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(phones)
}

// GetPhone returns a single phone by ID
//
// @Summary      Get a phone by ID
// @Tags         phones
// @Produce      json
// @Param        id   path      int  true  "Phone ID"
// @Success      200  {object}  model.Phone
// @Failure      400  {string}  string "invalid id"
// @Failure      404  {string}  string "not found"
// @Router       /phones/{id} [get]
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
//
// @Summary      Create a phone (admin only)
// @Tags         phones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body      model.Phone true "Phone data"
// @Success      201  {object}  model.Phone
// @Failure      400  {string}  string "invalid body"
// @Failure      401  {string}  string "unauthorized"
// @Failure      403  {string}  string "forbidden"
// @Router       /phones [post]
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

	if h.audit != nil {
		var uid *int
		if v, ok := r.Context().Value(middleware.UserIDKey).(int); ok {
			uid = &v
		}
		h.audit.Log(uid, "phone.created", "phone", strconv.Itoa(created.ID), clientIP(r),
			map[string]any{"brand": created.Brand, "model": created.Model, "price": created.Price})
	}

	// Return 201 Created
	w.WriteHeader(http.StatusCreated)

	// Return created phone
	json.NewEncoder(w).Encode(created)
}

// UpdatePhone updates an existing phone
//
// @Summary      Update a phone (admin only)
// @Tags         phones
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int         true "Phone ID"
// @Param        body body      model.Phone true "Updated phone data"
// @Success      200  {object}  model.Phone
// @Failure      400  {string}  string "invalid id or body"
// @Failure      401  {string}  string "unauthorized"
// @Failure      403  {string}  string "forbidden"
// @Failure      404  {string}  string "not found"
// @Router       /phones/{id} [put]
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

	if h.audit != nil {
		var uid *int
		if v, ok := r.Context().Value(middleware.UserIDKey).(int); ok {
			uid = &v
		}
		h.audit.Log(uid, "phone.updated", "phone", strconv.Itoa(p.ID), clientIP(r),
			map[string]any{"brand": p.Brand, "model": p.Model, "price": p.Price})
	}

	// Return updated phone
	json.NewEncoder(w).Encode(p)
}

// DeletePhone removes a phone by ID
//
// @Summary      Delete a phone (admin only)
// @Tags         phones
// @Security     BearerAuth
// @Param        id  path      int  true  "Phone ID"
// @Success      204 {string}  string ""
// @Failure      400 {string}  string "invalid id"
// @Failure      401 {string}  string "unauthorized"
// @Failure      403 {string}  string "forbidden"
// @Failure      404 {string}  string "not found"
// @Router       /phones/{id} [delete]
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

	if h.audit != nil {
		var uid *int
		if v, ok := r.Context().Value(middleware.UserIDKey).(int); ok {
			uid = &v
		}
		h.audit.Log(uid, "phone.deleted", "phone", strconv.Itoa(id), clientIP(r),
			map[string]any{"phone_id": id})
	}

	// Return 204 No Content (successful deletion)
	w.WriteHeader(http.StatusNoContent)
}
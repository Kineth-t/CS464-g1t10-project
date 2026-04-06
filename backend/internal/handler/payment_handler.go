package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// PaymentHandler handles HTTP requests for payment processing
type PaymentHandler struct {
	service *service.PaymentService
	audit   *service.AuditService
}

// Constructor
func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// SetAudit attaches an audit service for event logging.
func (h *PaymentHandler) SetAudit(s *service.AuditService) { h.audit = s }

// Pay processes a Stripe payment for the user's active cart
func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT context set by RequireAuth middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Decode the payment method ID from the request body
	var req service.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// Process the payment
	result, err := h.service.ProcessPayment(userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if h.audit != nil {
		h.audit.Log(&userID, "payment.completed", "order", result.PaymentID, clientIP(r),
			map[string]any{"order_id": result.PaymentID, "total": result.TotalAmount})
	}

	// Return the payment result
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
package handler

import (
	"encoding/json"
	"net/http"
	"log/slog"

	"github.com/Kineth-t/CS464-g1t10-project/internal/middleware"
	"github.com/Kineth-t/CS464-g1t10-project/internal/service"
)

// PaymentHandler handles HTTP requests for payment processing
type PaymentHandler struct {
	service *service.PaymentService
}

// Constructor
func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

// Pay processes a Stripe payment for the user's active cart
func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from JWT context set by RequireAuth middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Decode the payment method ID from the request body
	var req service.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// LOG: Request payload was malformed
        slog.Warn("payment request decode failed", 
            "user_id", userID, 
            "error", err.Error(),
        )

		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// LOG: The intent to pay (Crucial for tracking abandoned checkouts)
    slog.Info("processing payment start", 
        "user_id", userID, 
        "payment_method_id", req.PaymentMethodID,
    )

	// Process the payment
	result, err := h.service.ProcessPayment(userID, req)
	if err != nil {
		slog.Error("payment failed", "user_id", userID, "error", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Info("payment successful", "user_id", userID, "amount", result.TotalAmount)

	// Return the payment result
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}
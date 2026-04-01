package service

import (
	"errors"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// PaymentService handles Stripe payment processing and order creation
type PaymentService struct {
	cartRepo  repository.CartRepo  // Interface to access cart data
	phoneRepo repository.PhoneRepo // Interface to access phone data
	orderRepo repository.OrderRepo // Interface to persist orders
}

// Constructor — sets the Stripe secret key from environment on startup
func NewPaymentService(
	cartRepo repository.CartRepo,
	phoneRepo repository.PhoneRepo,
	orderRepo repository.OrderRepo,
) *PaymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &PaymentService{
		cartRepo:  cartRepo,
		phoneRepo: phoneRepo,
		orderRepo: orderRepo,
	}
}

// PaymentRequest holds the Stripe payment method ID from the frontend
type PaymentRequest struct {
	PaymentMethodID string `json:"payment_method_id"`
}

// PaymentResult is returned to the handler after a successful payment
type PaymentResult struct {
	PaymentID   string      `json:"payment_id"`   // Stripe pi_xxx — also used as order ID
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`       // "succeeded"
	Message     string      `json:"message"`
	Order       model.Order `json:"order"`        // Full order record
}

// ProcessPayment creates a Stripe PaymentIntent, checks out the cart,
// and persists the order using the Stripe payment ID as the order ID
func (s *PaymentService) ProcessPayment(userID int, req PaymentRequest) (PaymentResult, error) {
	// Ensure a payment method ID was provided
	if req.PaymentMethodID == "" {
		return PaymentResult{}, errors.New("payment_method_id is required")
	}

	// Get the user's active cart
	cart, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return PaymentResult{}, errors.New("no active cart")
	}

	// Ensure cart is not empty
	if len(cart.Items) == 0 {
		return PaymentResult{}, errors.New("cart is empty")
	}

	// Calculate total in cents — Stripe always works in the smallest currency unit
	totalCents := calculateTotalCents(cart.Items)

	// Get currency from environment — defaults to sgd
	currency := os.Getenv("STRIPE_CURRENCY")
	if currency == "" {
		currency = "sgd"
	}

	// Create and confirm a PaymentIntent in a single Stripe API call
	params := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(totalCents),
		Currency:      stripe.String(currency),
		PaymentMethod: stripe.String(req.PaymentMethodID),
		Confirm:       stripe.Bool(true),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"),
		},
	}

	// Call the Stripe API — real test charge happens here
	pi, err := paymentintent.New(params)
	if err != nil {
		// Stripe returns structured errors — extract the human-readable message
		if stripeErr, ok := err.(*stripe.Error); ok {
			return PaymentResult{}, errors.New(stripeErr.Msg)
		}
		return PaymentResult{}, errors.New("payment failed")
	}

	// Verify the intent actually succeeded
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return PaymentResult{}, errors.New("payment not completed: " + string(pi.Status))
	}

	// Checkout cart — deducts stock atomically inside a DB transaction
	if err := s.cartRepo.CheckoutCart(cart.ID, cart.Items); err != nil {
		// Payment succeeded but checkout failed — in production issue a Stripe refund here
		return PaymentResult{}, errors.New("payment succeeded but checkout failed: " + err.Error())
	}

	// Build order items from cart items
	orderItems := make([]model.CartItem, len(cart.Items))
	for i, item := range cart.Items {
		orderItems[i] = model.CartItem{
			PhoneID:  item.PhoneID,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	// Persist the order — use the Stripe payment intent ID as the order ID
	order, err := s.orderRepo.Create(model.Order{
		ID:     pi.ID,
		UserID: userID,
		Status: string(pi.Status),
		Total:  float64(totalCents) / 100,
		Items:  orderItems,
	})
	if err != nil {
		// Log it properly — don't silently swallow it
		// In production you'd also issue a Stripe refund here
		return PaymentResult{}, fmt.Errorf("order save failed: %w", err)
	}

	return PaymentResult{
		PaymentID:   pi.ID,
		TotalAmount: float64(totalCents) / 100,
		Status:      string(pi.Status),
		Message:     "payment successful",
		Order:       order,
	}, nil
}

// GetOrders returns all orders for a user
func (s *PaymentService) GetOrders(userID int) ([]model.Order, error) {
	return s.orderRepo.GetByUserID(userID)
}

// GetOrder returns a single order — verifies it belongs to the requesting user
func (s *PaymentService) GetOrder(userID int, orderID string) (model.Order, error) {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return model.Order{}, err
	}
	// Ensure the order belongs to the requesting user
	if order.UserID != userID {
		return model.Order{}, errors.New("order not found")
	}
	return order, nil
}

// calculateTotalCents converts the cart total to cents for Stripe
func calculateTotalCents(items []model.CartItem) int64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	// Multiply by 100 to convert to cents
	return int64(total * 100)
}
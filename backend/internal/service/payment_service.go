package service

import (
	"errors"
	"os"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// PaymentService handles Stripe payment processing
type PaymentService struct {
	cartRepo  repository.CartRepo  // Interface to access cart data
	phoneRepo repository.PhoneRepo // Interface to access phone data
}

// Constructor — sets the Stripe secret key from environment on startup
func NewPaymentService(cartRepo repository.CartRepo, phoneRepo repository.PhoneRepo) *PaymentService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &PaymentService{cartRepo: cartRepo, phoneRepo: phoneRepo}
}

// PaymentRequest holds the Stripe payment method ID created by the frontend.
// The frontend uses Stripe.js to tokenise card details
// The backend never sees raw card numbers, only this payment method ID.
// For backend testing, use Stripe's predefined test IDs like "pm_card_visa".
type PaymentRequest struct {
	PaymentMethodID string `json:"payment_method_id"`
}

// PaymentResult is returned to the handler after a successful payment
type PaymentResult struct {
	PaymentID   string  `json:"payment_id"`   // Stripe payment intent ID e.g. pi_xxx
	TotalAmount float64 `json:"total_amount"` // Total charged in dollars
	Status      string  `json:"status"`       // "succeeded"
	Message     string  `json:"message"`
}

// ProcessPayment creates a Stripe PaymentIntent and confirms it immediately.
// Validates the cart, calls Stripe, and checks out the cart atomically on success.
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

	// Create and confirm a PaymentIntent in a single API call
	params := &stripe.PaymentIntentParams{
		Amount:             stripe.Int64(totalCents),
		Currency:           stripe.String(currency),
		PaymentMethod:      stripe.String(req.PaymentMethodID),
		ConfirmationMethod: stripe.String("manual"),
		Confirm:            stripe.Bool(true),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:         stripe.Bool(true),
			AllowRedirects:  stripe.String("never"),
		},
	}

	// Call the Stripe API — this is where the real test charge happens
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

	// Checkout cart — deducts stock atomically inside a DB transaction.
	// If this fails after payment succeeded, in production you would
	// issue a refund via the Stripe refund API.
	if err := s.cartRepo.CheckoutCart(cart.ID, cart.Items); err != nil {
		return PaymentResult{}, errors.New("payment succeeded but checkout failed: " + err.Error())
	}

	return PaymentResult{
		PaymentID:   pi.ID,
		TotalAmount: float64(totalCents) / 100,
		Status:      string(pi.Status),
		Message:     "payment successful",
	}, nil
}

// calculateTotalCents converts the cart total to cents for Stripe
func calculateTotalCents(items []model.CartItem) int64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	// Multiply by 100 to convert dollars to cents
	return int64(total * 100)
}
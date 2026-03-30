package service

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// CartService contains the business logic for cart operations
type CartService struct {
	cartRepo  repository.CartRepo  // Interface to access cart data
	phoneRepo repository.PhoneRepo // Interface to access phone data
}

// Constructor
func NewCartService(cartRepo repository.CartRepo, phoneRepo repository.PhoneRepo) *CartService {
	return &CartService{cartRepo: cartRepo, phoneRepo: phoneRepo}
}

// GetCart retrieves the active cart for a given user
func (s *CartService) GetCart(userID int) (model.Cart, error) {
	return s.cartRepo.GetCartByUser(userID) // Returns error if no active cart exists
}

// AddToCart adds a phone item to the user's active cart.
// Uses a locked stock check so two users cannot both add the last item simultaneously.
func (s *CartService) AddToCart(userID, phoneID, quantity int) (model.CartItem, error) {
	// Lock the phone row and validate stock atomically — prevents race conditions
	price, err := s.phoneRepo.CheckStockAndReserve(phoneID, quantity)
	if err != nil {
		return model.CartItem{}, err // "phone not found" or "insufficient stock"
	}

	// Retrieve or create an active cart for the user
	cart := s.cartRepo.GetOrCreateActiveCart(userID)

	// Create a new cart item — price is captured at time of add
	item := model.CartItem{
		CartID:   cart.ID,
		PhoneID:  phoneID,
		Quantity: quantity,
		Price:    price,
	}
	// Add item to cart repository
	return s.cartRepo.AddItem(item), nil
}

// RemoveFromCart removes an item from the user's active cart
func (s *CartService) RemoveFromCart(userID, itemID int) error {
	// Get the user's active cart
	cart, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return errors.New("no active cart")
	}

	// Remove the item from the cart
	return s.cartRepo.RemoveItem(itemID, cart.ID)
}

// Checkout finalizes the user's cart, deducts stock, and marks it as checked out.
// All stock deductions happen inside a single database transaction with row-level locks
// so concurrent checkouts cannot both succeed against the same stock.
func (s *CartService) Checkout(userID int) error {
	// Get the user's active cart
	cart, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return errors.New("no active cart")
	}

	// Ensure cart is not empty
	if len(cart.Items) == 0 {
		return errors.New("cart is empty")
	}

	// Pass items into CheckoutCart — locking, validation, and deduction all happen
	// inside a single transaction in the repository layer
	return s.cartRepo.CheckoutCart(cart.ID, cart.Items)
}
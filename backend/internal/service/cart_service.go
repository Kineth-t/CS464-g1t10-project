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

// AddToCart adds a phone item to the user's active cart
func (s *CartService) AddToCart(userID, phoneID, quantity int) (model.CartItem, error) {
	// Check if the phone exists
	phone, err := s.phoneRepo.GetByID(phoneID)
	if err != nil {
		return model.CartItem{}, errors.New("phone not found")
	}

	// Ensure requested quantity is available
	if phone.Stock < quantity {
		return model.CartItem{}, errors.New("insufficient stock")
	}

	// Retrieve or create an active cart for the user
	cart := s.cartRepo.GetOrCreateActiveCart(userID)

	// Create a new cart item
	item := model.CartItem{
		CartID:   cart.ID,
		PhoneID:  phoneID,
		Quantity: quantity,
		Price:    phone.Price,
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

// Checkout finalizes the user's cart, deducts stock, and marks it as checked out
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

	// First pass: validate all items have sufficient stock before deducting anything
	// This prevents a partial checkout where some phones get deducted but a later one fails
	for _, item := range cart.Items {
		phone, err := s.phoneRepo.GetByID(item.PhoneID)
		if err != nil {
			return errors.New("a phone in your cart no longer exists")
		}
		if phone.Stock < item.Quantity {
			return errors.New("insufficient stock for one or more items in your cart")
		}
	}

	// Second pass: deduct stock for each item now that all items are confirmed available
	for _, item := range cart.Items {
		// Re-fetch phone to get the latest stock value
		phone, err := s.phoneRepo.GetByID(item.PhoneID)
		if err != nil {
			return err
		}

		// Deduct the purchased quantity from stock
		phone.Stock -= item.Quantity

		// Persist the updated stock
		if err := s.phoneRepo.Update(phone); err != nil {
			return err
		}
	}

	// Mark the cart as checked out in the repository
	return s.cartRepo.CheckoutCart(cart.ID)
}
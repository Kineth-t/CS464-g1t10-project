package service

import (
	"testing"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

func TestAddToCart_PhoneNotFound(t *testing.T) {
	svc := NewCartService(newMockCartRepo(), newMockPhoneRepo())
	_, err := svc.AddToCart(1, 999, 1)
	if err == nil || err.Error() != "phone not found" {
		t.Fatalf("expected 'phone not found', got %v", err)
	}
}

func TestAddToCart_InsufficientStock(t *testing.T) {
	phoneRepo := newMockPhoneRepo()
	phone := phoneRepo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 1})
	svc := NewCartService(newMockCartRepo(), phoneRepo)

	_, err := svc.AddToCart(1, phone.ID, 5)
	if err == nil || err.Error() != "insufficient stock" {
		t.Fatalf("expected 'insufficient stock', got %v", err)
	}
}

func TestAddToCart_Success(t *testing.T) {
	phoneRepo := newMockPhoneRepo()
	phone := phoneRepo.Create(model.Phone{Brand: "Samsung", Price: 799, Stock: 10})
	svc := NewCartService(newMockCartRepo(), phoneRepo)

	item, err := svc.AddToCart(1, phone.ID, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", item.Quantity)
	}
	if item.Price != 799 {
		t.Errorf("expected price 799, got %f", item.Price)
	}
	if item.PhoneID != phone.ID {
		t.Errorf("expected phone ID %d, got %d", phone.ID, item.PhoneID)
	}
}

func TestGetCart_NoCart(t *testing.T) {
	svc := NewCartService(newMockCartRepo(), newMockPhoneRepo())
	_, err := svc.GetCart(999)
	if err == nil {
		t.Fatal("expected error when no active cart exists")
	}
}

func TestGetCart_ExistingCart(t *testing.T) {
	cartRepo := newMockCartRepo()
	cartRepo.GetOrCreateActiveCart(1)
	svc := NewCartService(cartRepo, newMockPhoneRepo())

	cart, err := svc.GetCart(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cart.UserID != 1 {
		t.Errorf("expected userID 1, got %d", cart.UserID)
	}
	if cart.Status != "active" {
		t.Errorf("expected status active, got %s", cart.Status)
	}
}

func TestRemoveFromCart_NoCart(t *testing.T) {
	svc := NewCartService(newMockCartRepo(), newMockPhoneRepo())
	err := svc.RemoveFromCart(1, 1)
	if err == nil || err.Error() != "no active cart" {
		t.Fatalf("expected 'no active cart', got %v", err)
	}
}

func TestRemoveFromCart_Success(t *testing.T) {
	cartRepo := newMockCartRepo()
	phoneRepo := newMockPhoneRepo()
	phone := phoneRepo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 10})
	svc := NewCartService(cartRepo, phoneRepo)

	item, err := svc.AddToCart(1, phone.ID, 1)
	if err != nil {
		t.Fatalf("unexpected error adding to cart: %v", err)
	}
	if err := svc.RemoveFromCart(1, item.ID); err != nil {
		t.Fatalf("unexpected error removing from cart: %v", err)
	}
}

func TestCheckout_NoCart(t *testing.T) {
	svc := NewCartService(newMockCartRepo(), newMockPhoneRepo())
	err := svc.Checkout(1)
	if err == nil || err.Error() != "no active cart" {
		t.Fatalf("expected 'no active cart', got %v", err)
	}
}

func TestCheckout_EmptyCart(t *testing.T) {
	cartRepo := newMockCartRepo()
	cartRepo.GetOrCreateActiveCart(1)
	svc := NewCartService(cartRepo, newMockPhoneRepo())

	err := svc.Checkout(1)
	if err == nil || err.Error() != "cart is empty" {
		t.Fatalf("expected 'cart is empty', got %v", err)
	}
}

func TestCheckout_Success(t *testing.T) {
	cartRepo := newMockCartRepo()
	phoneRepo := newMockPhoneRepo()
	phone := phoneRepo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 10})
	svc := NewCartService(cartRepo, phoneRepo)

	if _, err := svc.AddToCart(1, phone.ID, 1); err != nil {
		t.Fatalf("unexpected error adding to cart: %v", err)
	}
	if err := svc.Checkout(1); err != nil {
		t.Fatalf("unexpected error during checkout: %v", err)
	}
}

func TestCheckout_CartNoLongerActiveAfterCheckout(t *testing.T) {
	cartRepo := newMockCartRepo()
	phoneRepo := newMockPhoneRepo()
	phone := phoneRepo.Create(model.Phone{Brand: "Apple", Price: 999, Stock: 10})
	svc := NewCartService(cartRepo, phoneRepo)

	if _, err := svc.AddToCart(1, phone.ID, 1); err != nil {
		t.Fatalf("unexpected error adding to cart: %v", err)
	}
	if err := svc.Checkout(1); err != nil {
		t.Fatalf("unexpected error during checkout: %v", err)
	}
	// After checkout, the active cart should be gone
	if _, err := svc.GetCart(1); err == nil {
		t.Fatal("expected no active cart after checkout")
	}
}

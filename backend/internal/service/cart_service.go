package service

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

type CartService struct {
	cartRepo  repository.CartRepo
	phoneRepo repository.PhoneRepo
}

func NewCartService(cartRepo repository.CartRepo, phoneRepo repository.PhoneRepo) *CartService {
	return &CartService{cartRepo: cartRepo, phoneRepo: phoneRepo}
}

func (s *CartService) GetCart(userID int) (model.Cart, error) {
	return s.cartRepo.GetCartByUser(userID)
}

func (s *CartService) AddToCart(userID, phoneID, quantity int) (model.CartItem, error) {
	phone, err := s.phoneRepo.GetByID(phoneID)
	if err != nil {
		return model.CartItem{}, errors.New("phone not found")
	}
	if phone.Stock < quantity {
		return model.CartItem{}, errors.New("insufficient stock")
	}
	cart := s.cartRepo.GetOrCreateActiveCart(userID)
	item := model.CartItem{
		CartID:   cart.ID,
		PhoneID:  phoneID,
		Quantity: quantity,
		Price:    phone.Price,
	}
	return s.cartRepo.AddItem(item), nil
}

func (s *CartService) RemoveFromCart(userID, itemID int) error {
	cart, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return errors.New("no active cart")
	}
	return s.cartRepo.RemoveItem(itemID, cart.ID)
}

func (s *CartService) Checkout(userID int) error {
	cart, err := s.cartRepo.GetCartByUser(userID)
	if err != nil {
		return errors.New("no active cart")
	}
	if len(cart.Items) == 0 {
		return errors.New("cart is empty")
	}
	return s.cartRepo.CheckoutCart(cart.ID)
}
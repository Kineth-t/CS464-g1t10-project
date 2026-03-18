package repository

import (
	"errors"
	"sync"
	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type CartRepository struct {
	mu       sync.RWMutex
	carts    map[int]model.Cart
	items    map[int]model.CartItem
	nextCart int
	nextItem int
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts:    make(map[int]model.Cart),
		items:    make(map[int]model.CartItem),
		nextCart: 1,
		nextItem: 1,
	}
}

func (r *CartRepository) GetOrCreateActiveCart(userID int) model.Cart {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, c := range r.carts {
		if c.UserID == userID && c.Status == "active" {
			return r.populateItems(c)
		}
	}
	cart := model.Cart{
		ID:     r.nextCart,
		UserID: userID,
		Status: "active",
		Items:  []model.CartItem{},
	}
	r.nextCart++
	r.carts[cart.ID] = cart
	return cart
}

func (r *CartRepository) GetCartByUser(userID int) (model.Cart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, c := range r.carts {
		if c.UserID == userID && c.Status == "active" {
			return r.populateItems(c), nil
		}
	}
	return model.Cart{}, errors.New("no active cart")
}

func (r *CartRepository) AddItem(item model.CartItem) model.CartItem {
	r.mu.Lock()
	defer r.mu.Unlock()
	item.ID = r.nextItem
	r.nextItem++
	r.items[item.ID] = item
	return item
}

func (r *CartRepository) RemoveItem(itemID, cartID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	item, ok := r.items[itemID]
	if !ok {
		return errors.New("cart item not found")
	}
	if item.CartID != cartID {
		return errors.New("item does not belong to this cart")
	}
	delete(r.items, itemID)
	return nil
}

func (r *CartRepository) CheckoutCart(cartID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cart, ok := r.carts[cartID]
	if !ok {
		return errors.New("cart not found")
	}
	cart.Status = "checked_out"
	r.carts[cartID] = cart
	return nil
}

func (r *CartRepository) populateItems(cart model.Cart) model.Cart {
	cart.Items = []model.CartItem{}
	for _, item := range r.items {
		if item.CartID == cart.ID {
			cart.Items = append(cart.Items, item)
		}
	}
	return cart
}
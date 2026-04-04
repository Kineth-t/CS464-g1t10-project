package handler

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// ── mockPhoneRepo ────────────────────────────────────────────────────────────

type mockPhoneRepo struct {
	phones map[int]model.Phone
	nextID int
}

func newMockPhoneRepo() *mockPhoneRepo {
	return &mockPhoneRepo{phones: make(map[int]model.Phone), nextID: 1}
}

func (m *mockPhoneRepo) GetAll() ([]model.Phone, error) {
	out := make([]model.Phone, 0, len(m.phones))
	for _, p := range m.phones {
		out = append(out, p)
	}
	return out, nil
}

func (m *mockPhoneRepo) GetByID(id int) (model.Phone, error) {
	p, ok := m.phones[id]
	if !ok {
		return model.Phone{}, errors.New("phone not found")
	}
	return p, nil
}

func (m *mockPhoneRepo) CheckStockAndReserve(phoneID, quantity int) (float64, error) {
	p, ok := m.phones[phoneID]
	if !ok {
		return 0, errors.New("phone not found")
	}
	if p.Stock < quantity {
		return 0, errors.New("insufficient stock")
	}
	updated := p
	updated.Stock -= quantity
	m.phones[phoneID] = updated
	return p.Price, nil
}

func (m *mockPhoneRepo) Create(p model.Phone) model.Phone {
	p.ID = m.nextID
	m.nextID++
	m.phones[p.ID] = p
	return p
}

func (m *mockPhoneRepo) Update(p model.Phone) error {
	if _, ok := m.phones[p.ID]; !ok {
		return errors.New("phone not found")
	}
	m.phones[p.ID] = p
	return nil
}

func (m *mockPhoneRepo) Delete(id int) error {
	if _, ok := m.phones[id]; !ok {
		return errors.New("phone not found")
	}
	delete(m.phones, id)
	return nil
}

// ── mockUserRepo ─────────────────────────────────────────────────────────────

type mockUserRepo struct {
	users  map[string]model.User
	nextID int
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]model.User), nextID: 1}
}

func (m *mockUserRepo) Create(u model.User) (model.User, error) {
	if _, exists := m.users[u.Username]; exists {
		return model.User{}, errors.New("username already taken")
	}
	u.ID = m.nextID
	m.nextID++
	m.users[u.Username] = u
	return u, nil
}

func (m *mockUserRepo) FindByUsername(username string) (model.User, error) {
	u, ok := m.users[username]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}

// ── mockCartRepo ─────────────────────────────────────────────────────────────

type mockCartRepo struct {
	carts    map[int]model.Cart
	items    map[int]model.CartItem
	nextCart int
	nextItem int
}

func newMockCartRepo() *mockCartRepo {
	return &mockCartRepo{
		carts:    make(map[int]model.Cart),
		items:    make(map[int]model.CartItem),
		nextCart: 1,
		nextItem: 1,
	}
}

func (m *mockCartRepo) GetOrCreateActiveCart(userID int) model.Cart {
	for _, c := range m.carts {
		if c.UserID == userID && c.Status == "active" {
			return m.withItems(c)
		}
	}
	cart := model.Cart{ID: m.nextCart, UserID: userID, Status: "active", Items: []model.CartItem{}}
	m.nextCart++
	m.carts[cart.ID] = cart
	return cart
}

func (m *mockCartRepo) GetCartByUser(userID int) (model.Cart, error) {
	for _, c := range m.carts {
		if c.UserID == userID && c.Status == "active" {
			return m.withItems(c), nil
		}
	}
	return model.Cart{}, errors.New("no active cart")
}

func (m *mockCartRepo) AddItem(item model.CartItem) model.CartItem {
	item.ID = m.nextItem
	m.nextItem++
	m.items[item.ID] = item
	return item
}

func (m *mockCartRepo) RemoveItem(itemID, cartID int) error {
	item, ok := m.items[itemID]
	if !ok {
		return errors.New("cart item not found")
	}
	if item.CartID != cartID {
		return errors.New("item does not belong to this cart")
	}
	delete(m.items, itemID)
	return nil
}

func (m *mockCartRepo) CheckoutCart(cartID int, items []model.CartItem) error {
	cart, ok := m.carts[cartID]
	if !ok {
		return errors.New("cart not found")
	}
	cart.Status = "checked_out"
	m.carts[cartID] = cart
	return nil
}

func (m *mockCartRepo) withItems(cart model.Cart) model.Cart {
	cart.Items = []model.CartItem{}
	for _, item := range m.items {
		if item.CartID == cart.ID {
			cart.Items = append(cart.Items, item)
		}
	}
	return cart
}

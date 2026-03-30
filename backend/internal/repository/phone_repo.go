package repository

import (
	"errors"
	"sync"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type PhoneRepo interface {
	GetAll() []model.Phone
	GetByID(id int) (model.Phone, error)
	CheckStockAndReserve(phoneID, quantity int) (float64, error) // locked stock check
	Create(p model.Phone) model.Phone
	Update(p model.Phone) error
	Delete(id int) error
}

type PhoneRepository struct {
	mu     sync.RWMutex
	store  map[int]model.Phone
	nextID int
}

func NewPhoneRepository() *PhoneRepository {
	return &PhoneRepository{store: make(map[int]model.Phone), nextID: 1}
}

func (r *PhoneRepository) GetAll() []model.Phone {
	r.mu.RLock()
	defer r.mu.RUnlock()
	phones := make([]model.Phone, 0, len(r.store))
	for _, p := range r.store {
		phones = append(phones, p)
	}
	return phones
}

func (r *PhoneRepository) GetByID(id int) (model.Phone, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.store[id]
	if !ok {
		return model.Phone{}, errors.New("phone not found")
	}
	return p, nil
}

func (r *PhoneRepository) Create(p model.Phone) model.Phone {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.nextID++
	r.store[p.ID] = p
	return p
}

func (r *PhoneRepository) Update(p model.Phone) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[p.ID]; !ok {
		return errors.New("phone not found")
	}
	r.store[p.ID] = p
	return nil
}

func (r *PhoneRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.store[id]; !ok {
		return errors.New("phone not found")
	}
	delete(r.store, id)
	return nil
}
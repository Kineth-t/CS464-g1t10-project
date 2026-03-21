package repository

import (
	"errors"
	"sync"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type UserRepo interface {
	Create(u model.User) (model.User, error)
	FindByUsername(username string) (model.User, error)
}

type UserRepository struct {
	mu     sync.RWMutex
	store  map[int]model.User
	byName map[string]model.User
	nextID int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		store:  make(map[int]model.User),
		byName: make(map[string]model.User),
		nextID: 1,
	}
}

func (r *UserRepository) Create(u model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byName[u.Username]; exists {
		return model.User{}, errors.New("username already taken")
	}
	u.ID = r.nextID
	r.nextID++
	r.store[u.ID] = u
	r.byName[u.Username] = u
	return u, nil
}

func (r *UserRepository) FindByUsername(username string) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byName[username]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return u, nil
}
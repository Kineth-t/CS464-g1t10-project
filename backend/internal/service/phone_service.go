package service

import (
	"context"
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// PhoneService handles business logic for phone operations
type PhoneService struct {
	repo repository.PhoneRepo // Repository interface for phone data access
	cache *repository.PhoneCache // Redis Cache
}

// Constructor
func NewPhoneService(repo repository.PhoneRepo, cache *repository.PhoneCache) *PhoneService {
	return &PhoneService{
		repo: repo,
		cache: cache,
	}
}

// ListPhones returns all phones
func (s *PhoneService) GetAll() ([]model.Phone, error) {
	ctx := context.Background()

	// Try to get from Redis first
	phones, err := s.cache.GetList(ctx)
	if err == nil && phones != nil {
		return phones, nil // Cache Hit! [cite: 530]
	}
	
	// Directly get all the phones in the database in case of Cache Miss
	phones, err = s.repo.GetAll() 
    if err != nil {
        return nil, err 
    }

	// Save to Redis for future requests
	_ = s.cache.SetList(ctx, phones)

    return phones, nil
}

// GetPhone returns a single phone by ID
func (s *PhoneService) GetPhone(id int) (model.Phone, error) {
	ctx := context.Background()

	// Check Redis first
    cachedPhone, err := s.cache.GetByID(ctx, id)
    if err == nil && cachedPhone != nil {
        return *cachedPhone, nil
    }

	// Fallback to Postgres
    phone, err := s.repo.GetByID(id)
    if err != nil {
        return model.Phone{}, err
    }

	// Populate cache for next time
    _ = s.cache.SetByID(ctx, id, phone)

	return phone, nil
}

// UpdatePhone updates phone data and clears stale cache
func (s *PhoneService) UpdatePhone(p model.Phone) error {
	// Returns error if phone ID does not exist
	err := s.repo.Update(p)
	if err != nil {
		return err
	}

	ctx := context.Background()
	
	// Clear the full list cache
	s.cache.Clear(ctx)

	// Clear the specific ID cache (NEW)
	s.cache.ClearByID(ctx, p.ID)


	return nil 
}

// DeletePhone removes a phone by ID
func (s *PhoneService) DeletePhone(id int) error {
	// Returns error if phone ID does not exist
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Clear the full list cache
	s.cache.Clear(ctx)
	
	// Clear the specific ID cache (NEW) 
	s.cache.ClearByID(ctx, id)

	return nil
}

// CreatePhone validates and creates a new phone
func (s *PhoneService) CreatePhone(p model.Phone) (model.Phone, error) {
	// Business validation
	if p.Price <= 0 {
		return model.Phone{}, errors.New("price must be greater than zero")
	}
	if p.Stock < 0 {
		return model.Phone{}, errors.New("stock cannot be negative")
	}

	createdPhone := s.repo.Create(p)

	// Invalidate cache so the new phone appears in the list 
	s.cache.Clear(context.Background())

	// Delegate creation to repository
	return createdPhone, nil
}
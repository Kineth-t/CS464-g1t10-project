package service

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// PhoneService handles business logic for phone operations
type PhoneService struct {
	repo repository.PhoneRepo // Repository interface for phone data access
}

// Constructor
func NewPhoneService(repo repository.PhoneRepo) *PhoneService {
	return &PhoneService{repo: repo}
}

// ListPhones returns all phones
func (s *PhoneService) ListPhones() []model.Phone {
	return s.repo.GetAll() // Directly fetches all phones from repository
}

// GetPhone returns a single phone by ID
func (s *PhoneService) GetPhone(id int) (model.Phone, error) {
	return s.repo.GetByID(id) // Returns error if phone not found
}

// UpdatePhone updates an existing phone's details
func (s *PhoneService) UpdatePhone(p model.Phone) error {
	return s.repo.Update(p) // Returns error if phone ID does not exist
}

// DeletePhone removes a phone by ID
func (s *PhoneService) DeletePhone(id int) error {
	return s.repo.Delete(id) // Returns error if phone ID does not exist
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

	// Delegate creation to repository
	return s.repo.Create(p), nil
}
package service

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

type PhoneService struct {
	repo repository.PhoneRepo
}

func NewPhoneService(repo repository.PhoneRepo) *PhoneService {
	return &PhoneService{repo: repo}
}

func (s *PhoneService) ListPhones() []model.Phone            { return s.repo.GetAll() }
func (s *PhoneService) GetPhone(id int) (model.Phone, error) { return s.repo.GetByID(id) }
func (s *PhoneService) UpdatePhone(p model.Phone) error      { return s.repo.Update(p) }
func (s *PhoneService) DeletePhone(id int) error             { return s.repo.Delete(id) }

func (s *PhoneService) CreatePhone(p model.Phone) (model.Phone, error) {
	if p.Price <= 0 {
		return model.Phone{}, errors.New("price must be greater than zero")
	}
	if p.Stock < 0 {
		return model.Phone{}, errors.New("stock cannot be negative")
	}
	return s.repo.Create(p), nil
}

func (s *PhoneService) PurchasePhone(id, quantity int) error {
	phone, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if phone.Stock < quantity {
		return errors.New("insufficient stock")
	}
	phone.Stock -= quantity
	return s.repo.Update(phone)
}
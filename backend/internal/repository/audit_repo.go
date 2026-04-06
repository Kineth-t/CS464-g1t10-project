package repository

import "github.com/Kineth-t/CS464-g1t10-project/internal/model"


type PhoneRepo interface {
	GetAll() ([]model.Phone, error)
	GetByID(id int) (model.Phone, error)
	CheckStockAndReserve(phoneID, quantity int) (float64, error) // locked stock check
	Create(p model.Phone) model.Phone
	Update(p model.Phone) error
	Delete(id int) error
}

// ONLY ONE AuditRepo declaration allowed
type AuditRepo interface {
    Create(log model.AuditLog) error
}
package repository

import "github.com/Kineth-t/CS464-g1t10-project/internal/model"

// OrderRepo defines the persistence operations for orders
type OrderRepo interface {
	Create(order model.Order) (model.Order, error)
	GetByUserID(userID int) ([]model.Order, error)
	GetByID(orderID string) (model.Order, error)
}

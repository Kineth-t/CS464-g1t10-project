package service

import (
	"errors"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
	"github.com/Kineth-t/CS464-g1t10-project/internal/repository"
)

// OrderService handles order retrieval
type OrderService struct {
	orderRepo repository.OrderRepo
}

// Constructor
func NewOrderService(orderRepo repository.OrderRepo) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

// GetOrders returns all orders for a user
func (s *OrderService) GetOrders(userID int) ([]model.Order, error) {
	orders, err := s.orderRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	// Always return an empty slice rather than null so the frontend can
	// safely iterate without a nil check
	if orders == nil {
		return []model.Order{}, nil
	}
	return orders, nil
}

// GetOrder returns a single order — verifies it belongs to the requesting user
func (s *OrderService) GetOrder(userID int, orderID string) (model.Order, error) {
	if orderID == "" {
		return model.Order{}, errors.New("order id is required")
	}

	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return model.Order{}, err
	}

	// Ownership check — prevent users from reading each other's orders
	if order.UserID != userID {
		return model.Order{}, errors.New("order not found")
	}

	return order, nil
}
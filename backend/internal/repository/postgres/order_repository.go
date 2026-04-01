package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// OrderRepository handles all order-related DB operations
type OrderRepository struct {
	db *pgxpool.Pool // Connection pool to PostgreSQL
}

// Constructor
func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts an order and all its items inside a single transaction
func (r *OrderRepository) Create(order model.Order) (model.Order, error) {
	// Begin transaction — order + items must be inserted together
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return model.Order{}, err
	}
	defer tx.Rollback(context.Background())

	// Insert the order using the Stripe payment intent ID as primary key
	err = tx.QueryRow(context.Background(),
		`INSERT INTO orders (id, user_id, status, total)
		 VALUES ($1, $2, $3, $4)
		 RETURNING created_at`,
		order.ID, order.UserID, order.Status, order.Total,
	).Scan(&order.CreatedAt)
	if err != nil {
		return model.Order{}, err
	}

	// Insert each order item linked to the order
	for i, item := range order.Items {
		err = tx.QueryRow(context.Background(),
			`INSERT INTO order_items (order_id, phone_id, quantity, price)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			order.ID, item.PhoneID, item.Quantity, item.Price,
		).Scan(&order.Items[i].ID)
		if err != nil {
			return model.Order{}, err
		}
	}

	return order, tx.Commit(context.Background())
}

// GetByUserID returns all orders for a given user ordered by most recent first
func (r *OrderRepository) GetByUserID(userID int) ([]model.Order, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, user_id, status, total, created_at
		 FROM orders WHERE user_id=$1
		 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt)
		// Populate items for each order
		o.Items, _ = r.getOrderItems(o.ID)
		orders = append(orders, o)
	}
	return orders, nil
}

// GetByID returns a single order by its Stripe payment intent ID
func (r *OrderRepository) GetByID(orderID string) (model.Order, error) {
	var o model.Order
	err := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, status, total, created_at
		 FROM orders WHERE id=$1`, orderID,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.Total, &o.CreatedAt)
	if err != nil {
		return model.Order{}, errors.New("order not found")
	}
	o.Items, _ = r.getOrderItems(o.ID)
	return o, nil
}

// getOrderItems fetches all items belonging to a given order
func (r *OrderRepository) getOrderItems(orderID string) ([]model.CartItem, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, phone_id, quantity, price
		 FROM order_items WHERE order_id=$1`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		rows.Scan(&item.ID, &item.PhoneID, &item.Quantity, &item.Price)
		items = append(items, item)
	}
	return items, nil
}
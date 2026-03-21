package postgres

import (
	"context" // Used for DB operations (timeouts, cancellation)
	"errors"

	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL connection pool

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

// CartRepository handles cart operations using PostgreSQL
type CartRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// Constructor to create repository
func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateActiveCart retrieves an active cart or creates one if none exists
func (r *CartRepository) GetOrCreateActiveCart(userID int) model.Cart {
	var cart model.Cart

	// Try to find existing active cart
	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, status 
		 FROM carts 
		 WHERE user_id=$1 AND status='active' 
		 LIMIT 1`, userID)

	err := row.Scan(&cart.ID, &cart.UserID, &cart.Status)

	// If no cart found → create a new one
	if err != nil {
		r.db.QueryRow(context.Background(),
			`INSERT INTO carts (user_id, status) 
			 VALUES ($1, 'active') 
			 RETURNING id, user_id, status`,
			userID,
		).Scan(&cart.ID, &cart.UserID, &cart.Status)
	}

	// Load cart items
	cart.Items = r.getItems(cart.ID)

	return cart
}

// GetCartByUser returns the active cart for a user
func (r *CartRepository) GetCartByUser(userID int) (model.Cart, error) {
	var cart model.Cart

	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, status 
		 FROM carts 
		 WHERE user_id=$1 AND status='active' 
		 LIMIT 1`, userID)

	// If not found → return error
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.Status); err != nil {
		return model.Cart{}, errors.New("no active cart")
	}

	// Load items into cart
	cart.Items = r.getItems(cart.ID)

	return cart, nil
}

// AddItem inserts a new item into cart_items table
func (r *CartRepository) AddItem(item model.CartItem) model.CartItem {

	// Insert item and return generated ID
	r.db.QueryRow(context.Background(),
		`INSERT INTO cart_items (cart_id, phone_id, quantity, price)
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id`,
		item.CartID, item.PhoneID, item.Quantity, item.Price,
	).Scan(&item.ID)

	return item
}

// RemoveItem deletes an item from the cart
func (r *CartRepository) RemoveItem(itemID, cartID int) error {

	// Delete only if item belongs to this cart
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM cart_items 
		 WHERE id=$1 AND cart_id=$2`,
		itemID, cartID)

	if err != nil {
		return err
	}

	// If no rows deleted → item not found
	if result.RowsAffected() == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

// CheckoutCart marks a cart as checked out
func (r *CartRepository) CheckoutCart(cartID int) error {

	// Update cart status
	_, err := r.db.Exec(context.Background(),
		`UPDATE carts 
		 SET status='checked_out' 
		 WHERE id=$1`, cartID)

	return err
}

// getItems retrieves all items belonging to a cart
func (r *CartRepository) getItems(cartID int) []model.CartItem {

	rows, err := r.db.Query(context.Background(),
		`SELECT id, cart_id, phone_id, quantity, price 
		 FROM cart_items 
		 WHERE cart_id=$1`, cartID)

	if err != nil {
		// Return empty list if query fails
		return []model.CartItem{}
	}
	defer rows.Close()

	var items []model.CartItem

	// Iterate through rows
	for rows.Next() {
		var item model.CartItem

		// Scan row into struct
		rows.Scan(&item.ID, &item.CartID, &item.PhoneID, &item.Quantity, &item.Price)

		items = append(items, item)
	}

	return items
}
package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Kineth-t/CS464-g1t10-project/internal/model"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateActiveCart(userID int) model.Cart {
	var cart model.Cart
	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, status FROM carts WHERE user_id=$1 AND status='active' LIMIT 1`, userID)
	err := row.Scan(&cart.ID, &cart.UserID, &cart.Status)
	if err != nil {
		r.db.QueryRow(context.Background(),
			`INSERT INTO carts (user_id, status) VALUES ($1, 'active') RETURNING id, user_id, status`,
			userID,
		).Scan(&cart.ID, &cart.UserID, &cart.Status)
	}
	cart.Items = r.getItems(cart.ID)
	return cart
}

func (r *CartRepository) GetCartByUser(userID int) (model.Cart, error) {
	var cart model.Cart
	row := r.db.QueryRow(context.Background(),
		`SELECT id, user_id, status FROM carts WHERE user_id=$1 AND status='active' LIMIT 1`, userID)
	if err := row.Scan(&cart.ID, &cart.UserID, &cart.Status); err != nil {
		return model.Cart{}, errors.New("no active cart")
	}
	cart.Items = r.getItems(cart.ID)
	return cart, nil
}

func (r *CartRepository) AddItem(item model.CartItem) model.CartItem {
	r.db.QueryRow(context.Background(),
		`INSERT INTO cart_items (cart_id, phone_id, quantity, price)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		item.CartID, item.PhoneID, item.Quantity, item.Price,
	).Scan(&item.ID)
	return item
}

func (r *CartRepository) RemoveItem(itemID, cartID int) error {
	result, err := r.db.Exec(context.Background(),
		`DELETE FROM cart_items WHERE id=$1 AND cart_id=$2`, itemID, cartID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("cart item not found")
	}
	return nil
}

func (r *CartRepository) CheckoutCart(cartID int) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE carts SET status='checked_out' WHERE id=$1`, cartID)
	return err
}

func (r *CartRepository) getItems(cartID int) []model.CartItem {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, cart_id, phone_id, quantity, price FROM cart_items WHERE cart_id=$1`, cartID)
	if err != nil {
		return []model.CartItem{}
	}
	defer rows.Close()
	var items []model.CartItem
	for rows.Next() {
		var item model.CartItem
		rows.Scan(&item.ID, &item.CartID, &item.PhoneID, &item.Quantity, &item.Price)
		items = append(items, item)
	}
	return items
}
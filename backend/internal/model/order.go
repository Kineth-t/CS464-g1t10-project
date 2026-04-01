package model

import "time"

type Order struct {
	ID        string     `json:"id"`         // Stripe payment intent ID e.g. pi_xxx
	UserID    int        `json:"user_id"`
	Status    string     `json:"status"`     // "succeeded", "failed"
	Total     float64    `json:"total"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time     `json:"created_at"`
}
package model

type Cart struct {
	ID     int        `json:"id"`
	UserID int        `json:"user_id"`
	Status string     `json:"status"` // "active" or "checked_out"
	Items  []CartItem `json:"items"`
}

type CartItem struct {
	ID       int     `json:"id"`
	CartID   int     `json:"cart_id"`
	PhoneID  int     `json:"phone_id"`
	PhoneName string  `json:"phone_name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
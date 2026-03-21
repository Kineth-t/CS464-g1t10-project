package model

type Phone struct {
	ID          int     `json:"id"`
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Description string  `json:"description"`
}
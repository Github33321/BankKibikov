package models

import "time"

type TaxiOrder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Price     float64   `json:"price"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

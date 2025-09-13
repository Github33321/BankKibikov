package models

import "time"

type Loan struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Amount   float64   `json:"amount"`
	IssuedAt time.Time `json:"issued_at"`
}

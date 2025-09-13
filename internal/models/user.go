package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string         `json:"id"`
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	Email        string         `json:"email"`
	Role         string         `json:"role"`
	OTPCode      sql.NullString `json:"-"`
	OTPExpiresAt sql.NullTime   `json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
}

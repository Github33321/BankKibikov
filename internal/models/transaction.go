package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullableString struct {
	sql.NullString
}

func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

type Transaction struct {
	ID        string         `json:"id"`
	FromUser  NullableString `json:"from_user"`
	ToUser    NullableString `json:"to_user"`
	Amount    float64        `json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
}

package repository

import (
	"BankKibikov/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxiRepository struct {
	db *pgxpool.Pool
}

func NewTaxiRepository(db *pgxpool.Pool) *TaxiRepository {
	return &TaxiRepository{db: db}
}

func (r *TaxiRepository) Create(ctx context.Context, o *models.TaxiOrder) error {
	row := r.db.QueryRow(ctx,
		`INSERT INTO taxi_orders (user_id, from_address, to_address, price, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id, created_at`,
		o.UserID, o.From, o.To, o.Price, o.Status,
	)
	return row.Scan(&o.ID, &o.CreatedAt)
}

func (r *TaxiRepository) GetByID(ctx context.Context, id string) (*models.TaxiOrder, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, from_address, to_address, price, status, created_at
		 FROM taxi_orders WHERE id=$1`, id,
	)
	var o models.TaxiOrder
	if err := row.Scan(&o.ID, &o.UserID, &o.From, &o.To, &o.Price, &o.Status, &o.CreatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *TaxiRepository) GetByUser(ctx context.Context, userID string) ([]models.TaxiOrder, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, from_address, to_address, price, status, created_at
		 FROM taxi_orders WHERE user_id=$1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.TaxiOrder
	for rows.Next() {
		var o models.TaxiOrder
		if err := rows.Scan(&o.ID, &o.UserID, &o.From, &o.To, &o.Price, &o.Status, &o.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, nil
}

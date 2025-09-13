package repository

import (
	"BankKibikov/internal/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, t *models.Transaction) error {
	if !t.FromUser.Valid && !t.ToUser.Valid {
		return errors.New("invalid transaction: both from_user and to_user are NULL")
	}

	if !t.FromUser.Valid {
		row := r.db.QueryRow(ctx,
			"INSERT INTO transactions (from_user, to_user, amount) VALUES (NULL, $1, $2) RETURNING id, created_at",
			t.ToUser.String, t.Amount,
		)
		return row.Scan(&t.ID, &t.CreatedAt)
	}

	if !t.ToUser.Valid {
		row := r.db.QueryRow(ctx,
			"INSERT INTO transactions (from_user, to_user, amount) VALUES ($1, NULL, $2) RETURNING id, created_at",
			t.FromUser.String, t.Amount,
		)
		return row.Scan(&t.ID, &t.CreatedAt)
	}

	row := r.db.QueryRow(ctx,
		"INSERT INTO transactions (from_user, to_user, amount) VALUES ($1, $2, $3) RETURNING id, created_at",
		t.FromUser.String, t.ToUser.String, t.Amount,
	)
	return row.Scan(&t.ID, &t.CreatedAt)
}

func (r *TransactionRepository) GetByUser(ctx context.Context, userID string) ([]models.Transaction, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, from_user, to_user, amount, created_at FROM transactions "+
			"WHERE from_user=$1 OR to_user=$1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromUser, &t.ToUser, &t.Amount, &t.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, t)
	}
	return list, nil
}

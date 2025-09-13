package repository

import (
	"BankKibikov/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepository struct {
	db *pgxpool.Pool
}

func NewAccountRepository(db *pgxpool.Pool) *AccountRepository {
	return &AccountRepository{db: db}
}
func (r *AccountRepository) CreateAccount(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, "INSERT INTO accounts (user_id, balance) VALUES ($1, 0)", userID)
	return err
}

func (r *AccountRepository) GetByUserID(ctx context.Context, userID string) (*models.Account, error) {
	row := r.db.QueryRow(ctx, "SELECT id, user_id, balance, created_at FROM accounts WHERE user_id=$1", userID)
	var acc models.Account
	if err := row.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.CreatedAt); err != nil {
		return nil, err
	}
	return &acc, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, userID string, delta float64) error {
	_, err := r.db.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE user_id=$2", delta, userID)
	return err
}

func (r *AccountRepository) Deposit(ctx context.Context, userID string, amount float64) error {
	_, err := r.db.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE user_id=$2", amount, userID)
	return err
}
func (r *AccountRepository) Withdraw(ctx context.Context, userID string, amount float64) error {
	_, err := r.db.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE user_id=$2", amount, userID)
	return err
}

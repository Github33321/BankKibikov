package repository

import (
	"BankKibikov/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanRepository struct {
	db *pgxpool.Pool
}

func NewLoanRepository(db *pgxpool.Pool) *LoanRepository {
	return &LoanRepository{db: db}
}

func (r *LoanRepository) Create(ctx context.Context, userID string, amount float64) (*models.Loan, error) {
	row := r.db.QueryRow(ctx,
		"INSERT INTO loans (user_id, amount) VALUES ($1, $2) RETURNING id, user_id, amount, issued_at",
		userID, amount,
	)

	var loan models.Loan
	if err := row.Scan(&loan.ID, &loan.UserID, &loan.Amount, &loan.IssuedAt); err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *LoanRepository) GetByID(ctx context.Context, id string) (*models.Loan, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, user_id, amount, issued_at FROM loans WHERE id=$1", id,
	)

	var loan models.Loan
	if err := row.Scan(&loan.ID, &loan.UserID, &loan.Amount, &loan.IssuedAt); err != nil {
		return nil, err
	}
	return &loan, nil
}

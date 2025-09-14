package service

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
	"database/sql"
	"time"
)

type LoanService struct {
	repo        *repository.LoanRepository
	accountRepo *repository.AccountRepository
	txRepo      *repository.TransactionRepository
}

func NewLoanService(
	repo *repository.LoanRepository,
	accountRepo *repository.AccountRepository,
	txRepo *repository.TransactionRepository,
) *LoanService {
	return &LoanService{repo: repo, accountRepo: accountRepo, txRepo: txRepo}
}

func (s *LoanService) CreateLoan(ctx context.Context, userID string, amount float64) (*models.Loan, error) {
	// создаём запись в таблице loans
	loan, err := s.repo.Create(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	// пополняем счёт пользователя
	if err := s.accountRepo.Deposit(ctx, userID, amount); err != nil {
		return nil, err
	}

	// создаём запись в истории транзакций (loan -> user)
	tx := &models.Transaction{
		FromUser:  models.NullableString{sql.NullString{Valid: false}},                // NULL = займ от системы
		ToUser:    models.NullableString{sql.NullString{String: userID, Valid: true}}, // получатель
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, err
	}

	return loan, nil
}

func (s *LoanService) GetLoanWithInterest(ctx context.Context, id string) (float64, error) {
	loan, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}

	seconds := time.Since(loan.IssuedAt).Seconds()
	current := loan.Amount * pow(1.1, int(seconds))

	if current > 99999999 {
		current = 9999999
	}

	return current, nil
}

func pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

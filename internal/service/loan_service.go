package service

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
	"time"
)

type LoanService struct {
	repo        *repository.LoanRepository
	accountRepo *repository.AccountRepository
}

func NewLoanService(repo *repository.LoanRepository, accountRepo *repository.AccountRepository) *LoanService {
	return &LoanService{repo: repo, accountRepo: accountRepo}
}

func (s *LoanService) CreateLoan(ctx context.Context, userID string, amount float64) (*models.Loan, error) {
	loan, err := s.repo.Create(ctx, userID, amount)
	if err != nil {
		return nil, err
	}

	if err := s.accountRepo.Deposit(ctx, userID, amount); err != nil {
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

	return current, nil
}

func pow(base float64, exp int) float64 {
	result := 1.0
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}

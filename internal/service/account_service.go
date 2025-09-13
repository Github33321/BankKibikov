package service

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
	"database/sql"
	"errors"
)

type AccountService struct {
	accRepo *repository.AccountRepository
	txRepo  *repository.TransactionRepository
}

func NewAccountService(accRepo *repository.AccountRepository, txRepo *repository.TransactionRepository) *AccountService {
	return &AccountService{accRepo: accRepo, txRepo: txRepo}
}

// Получение баланса
func (s *AccountService) GetBalance(ctx context.Context, userID string) (*models.Account, error) {
	return s.accRepo.GetByUserID(ctx, userID)
}

// Перевод денег
func (s *AccountService) Transfer(ctx context.Context, fromUser, toUser string, amount float64) error {
	acc, err := s.accRepo.GetByUserID(ctx, fromUser)
	if err != nil {
		return err
	}
	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}

	if err := s.accRepo.UpdateBalance(ctx, fromUser, -amount); err != nil {
		return err
	}
	if err := s.accRepo.UpdateBalance(ctx, toUser, amount); err != nil {
		return err
	}

	tx := &models.Transaction{
		FromUser: models.NullableString{sql.NullString{String: fromUser, Valid: true}},
		ToUser:   models.NullableString{sql.NullString{String: toUser, Valid: true}},
		Amount:   amount,
	}
	return s.txRepo.Create(ctx, tx)
}

// Депозит
func (s *AccountService) Deposit(ctx context.Context, userID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if err := s.accRepo.Deposit(ctx, userID, amount); err != nil {
		return err
	}

	tx := &models.Transaction{
		FromUser: models.NullableString{sql.NullString{Valid: false}},                // NULL
		ToUser:   models.NullableString{sql.NullString{String: userID, Valid: true}}, // получатель
		Amount:   amount,
	}
	return s.txRepo.Create(ctx, tx)
}

// Снятие (withdraw)
func (s *AccountService) Withdraw(ctx context.Context, userID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	acc, err := s.accRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}

	if err := s.accRepo.Withdraw(ctx, userID, amount); err != nil {
		return err
	}

	tx := &models.Transaction{
		FromUser: models.NullableString{sql.NullString{String: userID, Valid: true}}, // кто снимает
		ToUser:   models.NullableString{sql.NullString{Valid: false}},                // NULL
		Amount:   amount,
	}
	return s.txRepo.Create(ctx, tx)
}

// История операций
func (s *AccountService) GetTransactions(ctx context.Context, userID string) ([]models.Transaction, error) {
	return s.txRepo.GetByUser(ctx, userID)
}

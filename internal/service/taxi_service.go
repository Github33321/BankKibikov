package service

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
)

type TaxiService struct {
	repo           *repository.TaxiRepository
	accountService *AccountService
}

func NewTaxiService(repo *repository.TaxiRepository, accountService *AccountService) *TaxiService {
	return &TaxiService{repo: repo, accountService: accountService}
}

func (s *TaxiService) OrderTaxi(ctx context.Context, userID, from, to string, price float64) (*models.TaxiOrder, error) {
	if err := s.accountService.Withdraw(ctx, userID, price); err != nil {
		return nil, err
	}

	order := &models.TaxiOrder{
		UserID: userID,
		From:   from,
		To:     to,
		Price:  price,
		Status: "created",
	}
	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *TaxiService) GetOrder(ctx context.Context, id string) (*models.TaxiOrder, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaxiService) GetUserOrders(ctx context.Context, userID string) ([]models.TaxiOrder, error) {
	return s.repo.GetByUser(ctx, userID)
}

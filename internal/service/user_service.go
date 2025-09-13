package service

import (
	"BankKibikov/internal/models"
	"BankKibikov/internal/repository"
	"context"
)

type UserService struct {
	repo        *repository.UserRepository
	accountRepo *repository.AccountRepository
}

func NewUserService(userRepo *repository.UserRepository, accountRepo *repository.AccountRepository) *UserService {
	return &UserService{repo: userRepo, accountRepo: accountRepo}
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) error {
	if u.Password == "" {
		u.Password = "12345"
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}

	return s.accountRepo.CreateAccount(ctx, u.ID)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetAll(ctx)
}

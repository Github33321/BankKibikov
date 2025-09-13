package security

import (
	"BankKibikov/internal/repository"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

type AuthService struct {
	repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func generateOTP() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%06d", int(b[0])<<16|int(b[1])<<8|int(b[2])%1000000)
}

func (s *AuthService) RequestOTP(ctx context.Context, username, password string) error {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil || user.Password != password {
		return errors.New("invalid username or password")
	}

	otp := generateOTP()
	expires := time.Now().Add(5 * time.Minute)

	if err := s.repo.SaveOTP(ctx, user.ID, otp, expires); err != nil {
		return err
	}

	fmt.Printf("DEBUG OTP for user %s: %s (expires %s)\n", user.Username, otp, expires.Format(time.RFC3339))

	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, username, otp string) (string, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("user not found")
	}

	if !user.OTPCode.Valid || user.OTPCode.String != otp {
		return "", errors.New("invalid or expired OTP")
	}

	if !user.OTPExpiresAt.Valid || time.Now().After(user.OTPExpiresAt.Time) {
		return "", errors.New("invalid or expired OTP")
	}

	fmt.Println("DEBUG otp from DB =", "["+user.OTPCode.String+"]")
	fmt.Println("DEBUG otp from request =", "["+otp+"]")

	return "login successful", nil
}

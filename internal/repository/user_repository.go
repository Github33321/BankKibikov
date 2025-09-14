package repository

import (
	"BankKibikov/internal/models"
	"context"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	row := r.db.QueryRow(ctx,
		"INSERT INTO users (username, password, email, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id",
		u.Username, u.Password, u.Email,
	)
	return row.Scan(&u.ID)
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)", email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, username, password, created_at FROM users WHERE id = $1", id,
	)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, username, password, email, otp_code, otp_expires_at, created_at FROM users WHERE username = $1",
		username,
	)

	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password,
		&user.Email, &user.OTPCode, &user.OTPExpiresAt, &user.CreatedAt); err != nil {
		return nil, err
	}

	if user.OTPCode.Valid {
		user.OTPCode.String = strings.TrimSpace(user.OTPCode.String)
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, username, password, email, role, created_at  FROM users ORDER BY created_at`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) SaveOTP(ctx context.Context, userID, otp string, expires time.Time) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET otp_code=$1, otp_expires_at=$2 WHERE id=$3",
		otp, expires, userID,
	)
	return err
}

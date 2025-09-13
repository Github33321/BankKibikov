package repository

import (
	"BankKibikov/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsRepository struct {
	db *pgxpool.Pool
}

func NewNewsRepository(db *pgxpool.Pool) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) GetAll(ctx context.Context) ([]models.News, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, content, created_at FROM news ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []models.News
	for rows.Next() {
		var n models.News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.CreatedAt); err != nil {
			return nil, err
		}
		news = append(news, n)
	}
	return news, nil
}

func (r *NewsRepository) Create(ctx context.Context, n *models.News) error {
	row := r.db.QueryRow(ctx,
		"INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id, created_at",
		n.Title, n.Content,
	)
	return row.Scan(&n.ID, &n.CreatedAt)
}

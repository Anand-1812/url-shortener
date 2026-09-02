// Package repository for the pg repository
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRecords struct {
	ID          int64     `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		pool: pool,
	}
}

func (r *URLRepository) Create(ctx context.Context, originalURL, shortCode string) (URLRecords, error) {
	query := `
		INSERT INTO urls (original_url, short_url)
		VALUES ($1, $2)
		RETURNING id, original_url, short_url, clicks, created_at;
	`

	var records URLRecords
	err := r.pool.QueryRow(ctx, query, originalURL, shortCode).Scan(
		&records.ID,
		&records.OriginalURL,
		&records.ShortURL,
		&records.Clicks,
		&records.CreatedAt,
	)
	if err != nil {
		return URLRecords{}, fmt.Errorf("repository: failed to insert url: %w", err)
	}

	return records, nil
}

func (r *URLRepository) FindAndIncrementsClicks(ctx context.Context, shortCode string) (URLRecords, error) {
	query := `
		UPDATE urls
		SET clicks = clicks + 1
		where short_url = $1
		RETURNING id, original_url, short_url, clicks, created_at
	`

	var records URLRecords
	err := r.pool.QueryRow(ctx, query, shortCode).Scan(
		&records.ID,
		&records.OriginalURL,
		&records.ShortURL,
		&records.Clicks,
		&records.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return URLRecords{}, fmt.Errorf("repository: database error: %w", err)
		}
		return URLRecords{}, fmt.Errorf("repository: failed to fetch url: %w", err)
	}

	return records, nil
}

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

var ErrNotFound = errors.New("url record not found")

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

func (r *URLRepository) Ping(ctx context.Context) error {
	return r.pool.Ping(ctx)
}

func (r *URLRepository) CreateWithBase62(ctx context.Context, originalURL string, encodeFn func(int64) string) (URLRecords, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return URLRecords{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	insertQuery := `
		INSERT INTO urls (original_url, short_url)
		VALUES ($1, '')
		RETURNING id, clicks, created_at;
	`

	var record URLRecords
	record.OriginalURL = originalURL
	err = tx.QueryRow(ctx, insertQuery, originalURL).Scan(
		&record.ID,
		&record.Clicks,
		&record.CreatedAt,
	)
	if err != nil {
		return URLRecords{}, fmt.Errorf("repository: failed to insert initial row: %w", err)
	}

	const idOffset = 10000000
	record.ShortURL = encodeFn(record.ID + idOffset)

	// Step 3: Update row with the generated code
	updateQuery := `UPDATE urls SET short_url = $1 WHERE id = $2;`
	_, err = tx.Exec(ctx, updateQuery, record.ShortURL, record.ID)
	if err != nil {
		return URLRecords{}, fmt.Errorf("repository: failed to update short_url: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return URLRecords{}, fmt.Errorf("repository: failed to commit tx: %w", err)
	}

	return record, nil
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
			return URLRecords{}, ErrNotFound
		}
		return URLRecords{}, fmt.Errorf("repository: failed to fetch url: %w", err)
	}

	return records, nil
}

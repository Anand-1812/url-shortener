package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InsertURL(ctx context.Context, pool *pgxpool.Pool, originalURL, shortURL string) (int, error) {
	query := `
		INSERT INTO urls (original_url, short_url, clicks) 
		VALUES ($1, $2)
		RETURNING id;
	`

	var id int
	err := pool.QueryRow(ctx, query, originalURL, shortURL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unable to insert data: %w", err)
	}

	return id, nil
}

func GetOriginalURL(ctx context.Context, pool *pgxpool.Pool, shortURL string) (string, error) {
	query := `
		UPDATE urls
		SET clicks = clicks + 1
		WHERE short_url = $1
		RETURNING original_url
	`

	var originalURL string
	err := pool.QueryRow(ctx, query, shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("url not found")
		}

		return "", fmt.Errorf("query failed: %w", err)
	}

	return originalURL, nil
}

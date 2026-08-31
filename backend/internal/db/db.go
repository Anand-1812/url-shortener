// Package db give connection to the pg database
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context, dbConnString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dbConnString)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

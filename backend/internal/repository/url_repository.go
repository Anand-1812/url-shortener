package repository

import "github.com/jackc/pgx/v5/pgxpool"

type URLRepository struct {
	pool *pgxpool.Pool
}

func NewURLRepository(pool *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		pool: pool,
	}
}

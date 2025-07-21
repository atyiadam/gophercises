package repository

import (
	"context"
	"urlshortener/internal/db"
)

// URLRepository defines the interface for our URL data operations
type URLRepository interface {
	GetShortURL(ctx context.Context, shortPath string) (db.ShortUrl, error)
}

type PostgresRepository struct {
	queries *db.Queries
}

func NewPostgresRepository(dbtx db.DBTX) *PostgresRepository {
	return &PostgresRepository{
		queries: db.New(dbtx),
	}
}

func (r *PostgresRepository) GetShortURL(ctx context.Context, shortPath string) (db.ShortUrl, error) {
	return r.queries.GetShortURLByPath(ctx, shortPath)
}

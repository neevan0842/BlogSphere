package categories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type svc struct {
	repo *sqlc.Queries
	db   *pgxpool.Pool
}

func NewService(repo *sqlc.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) GetCategories(ctx context.Context) ([]sqlc.Category, error) {
	return s.repo.GetCategories(ctx)
}

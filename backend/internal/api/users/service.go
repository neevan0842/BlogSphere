package users

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
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

func (s *svc) getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by ID: %s", err.Error())
	}
	return user, nil
}

func (s *svc) getUserByUsername(ctx context.Context, username pgtype.Text) (sqlc.User, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by username: %s", err.Error())
	}
	return user, nil
}

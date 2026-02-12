package posts

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
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

func (s *svc) getPostsPaginated(ctx context.Context, search string, limit, offset int, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error) {
	// Fetch posts based on search query with pagination
	posts, err := s.repo.GetPostBySearchPaginated(ctx, sqlc.GetPostBySearchPaginatedParams{
		Column1: pgtype.Text{String: search, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return []common.PostCardDTO{}, fmt.Errorf("failed to fetch posts: %s", err.Error())
	}

	// Enrich posts with additional details
	return common.EnrichPostsWithDetails(ctx, s.repo, posts, requestingUserID)
}

func (s *svc) getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by ID: %s", err.Error())
	}
	return user, nil
}

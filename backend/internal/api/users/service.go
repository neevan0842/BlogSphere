package users

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

func (s *svc) updateUserDescription(ctx context.Context, userID pgtype.UUID, description pgtype.Text) (sqlc.User, error) {
	user, err := s.repo.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:          userID,
		Description: description,
	})
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to update user description: %s", err.Error())
	}
	return user, nil
}

func (s *svc) getPostsByUsername(ctx context.Context, username pgtype.Text, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error) {
	posts, err := s.repo.GetPostsByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by username: %s", err.Error())
	}

	return common.EnrichPostsWithDetails(ctx, s.repo, posts, requestingUserID)
}

func (s *svc) getLikedPostsByUsername(ctx context.Context, username pgtype.Text, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error) {
	posts, err := s.repo.GetPostsLikedByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get liked posts by username: %s", err.Error())
	}

	return common.EnrichPostsWithDetails(ctx, s.repo, posts, requestingUserID)
}

func (s *svc) deleteUserByID(ctx context.Context, userID pgtype.UUID) error {
	err := s.repo.DeleteUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %s", err.Error())
	}
	return nil
}

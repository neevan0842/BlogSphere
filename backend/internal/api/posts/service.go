package posts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

func (s *svc) getPostBySlug(ctx context.Context, slug string, requestingUserID *pgtype.UUID) (common.PostCardDTO, error) {
	post, err := s.repo.GetPostBySlug(ctx, slug)
	if err != nil {
		return common.PostCardDTO{}, fmt.Errorf("failed to get post by slug: %s", err.Error())
	}

	posts, err := common.EnrichPostsWithDetails(ctx, s.repo, []sqlc.Post{post}, requestingUserID)
	if err != nil {
		return common.PostCardDTO{}, fmt.Errorf("failed to enrich post details: %s", err.Error())
	}

	if len(posts) == 0 {
		return common.PostCardDTO{}, fmt.Errorf("post not found")
	}

	return posts[0], nil
}

func (s *svc) getCommentsByPostSlug(ctx context.Context, slug string) ([]common.CommentDTO, error) {
	comments, err := s.repo.GetCommentsByPostSlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post slug: %s", err.Error())
	}

	return common.EnrichCommentsWithAuthors(ctx, s.repo, comments)
}

func (s *svc) togglePostLike(ctx context.Context, postID pgtype.UUID, userID pgtype.UUID) (bool, error) {
	// Check if user has already liked the post
	_, err := s.repo.GetPostLike(ctx, sqlc.GetPostLikeParams{
		PostID: postID,
		UserID: userID,
	})

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("failed to check existing like: %s", err.Error())
	}

	if errors.Is(err, pgx.ErrNoRows) {
		// User has not liked the post, so add like
		_, err := s.repo.CreatePostLike(ctx, sqlc.CreatePostLikeParams{
			PostID: postID,
			UserID: userID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to like post: %s", err.Error())
		}
		return true, nil // Post is now liked
	} else {
		// User has already liked the post, so remove like
		err := s.repo.DeletePostLike(ctx, sqlc.DeletePostLikeParams{
			PostID: postID,
			UserID: userID,
		})
		if err != nil {
			return false, fmt.Errorf("failed to unlike post: %s", err.Error())
		}
		return false, nil // Post is now unliked
	}
}

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
	"github.com/neevan0842/BlogSphere/backend/utils"
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

func (s *svc) getPostsPaginated(ctx context.Context, search string, categorySlug string, limit, offset int, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error) {
	// Fetch posts based on search query with pagination
	posts, err := s.repo.GetPostBySearchAndCategoryPaginated(ctx, sqlc.GetPostBySearchAndCategoryPaginatedParams{
		CategorySlug: pgtype.Text{String: categorySlug, Valid: categorySlug != ""},
		Search:       pgtype.Text{String: search, Valid: search != ""},
		Limit:        int32(limit),
		Offset:       int32(offset),
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

func (s *svc) CreatePost(ctx context.Context, title string, body string, authorID string, categoryIDs []string) (common.PostCardDTO, error) {
	var createdPost sqlc.Post
	err := common.ExecTx(ctx, s.db, func(q *sqlc.Queries) error {
		// Generate slug from title
		slug := utils.GenerateSlug(title)
		authorUUID, err := utils.StrToUUID(authorID)
		if err != nil {
			return fmt.Errorf("invalid author ID: %s", err.Error())
		}

		// Create the post
		post, err := q.CreatePost(ctx, sqlc.CreatePostParams{
			Title:       title,
			Body:        body,
			Slug:        slug,
			AuthorID:    authorUUID,
			IsPublished: true,
		})
		if err != nil {
			return fmt.Errorf("failed to create post: %s", err.Error())
		}

		// Associate post with categories
		postIDUUIDs := make([]pgtype.UUID, len(categoryIDs))
		categoryIDsUUIDs := make([]pgtype.UUID, len(categoryIDs))
		for i, catID := range categoryIDs {
			postIDUUIDs[i] = post.ID
			categoryIDsUUIDs[i], err = utils.StrToUUID(catID)
			if err != nil {
				return fmt.Errorf("invalid category ID: %s", err.Error())
			}
		}

		// Batch insert post-category associations using 'unnest'
		err = q.BatchCreatePostCategories(ctx, sqlc.BatchCreatePostCategoriesParams{
			PostID:     postIDUUIDs,
			CategoryID: categoryIDsUUIDs})
		if err != nil {
			return fmt.Errorf("failed to associate post with category: %s", err.Error())
		}

		createdPost = post
		return nil
	})
	if err != nil {
		return common.PostCardDTO{}, err
	}
	posts, err := common.EnrichPostsWithDetails(ctx, s.repo, []sqlc.Post{createdPost}, &createdPost.AuthorID)
	if err != nil {
		return common.PostCardDTO{}, err
	}
	return posts[0], nil
}

func (s *svc) DeletePost(ctx context.Context, postID string, userID string) error {
	postIDUUID, err := utils.StrToUUID(postID)
	if err != nil {
		return fmt.Errorf("invalid post ID: %s", err.Error())
	}

	// Fetch the post to verify ownership
	post, err := s.repo.GetPostByID(ctx, postIDUUID)
	if err != nil {
		return fmt.Errorf("no post found with the given ID: %s", err.Error())
	}

	// Check if the requesting user is the author of the post
	if post.AuthorID.String() != userID {
		return fmt.Errorf("unauthorized: user does not own the post")
	}

	// Delete the post
	err = s.repo.DeletePost(ctx, postIDUUID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %s", err.Error())
	}

	return nil
}

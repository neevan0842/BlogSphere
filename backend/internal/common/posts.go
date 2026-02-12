package common

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"golang.org/x/sync/errgroup"
)

// EnrichPostsWithDetails fetches and attaches categories, likes, comments, and user-liked status to posts
func EnrichPostsWithDetails(ctx context.Context, repo *sqlc.Queries, posts []sqlc.Post, requestingUserID *pgtype.UUID) ([]PostCardDTO, error) {
	if len(posts) == 0 {
		return []PostCardDTO{}, nil
	}

	// extract post IDs and author IDs
	postIDs := make([]pgtype.UUID, len(posts))
	authorIDmap := make(map[string]bool)

	for i, post := range posts {
		postIDs[i] = post.ID
		authorIDmap[post.AuthorID.String()] = true
	}

	authorIDs := make([]pgtype.UUID, 0, len(authorIDmap))
	for authorID := range authorIDmap {
		id, _ := utils.StrToUUID(authorID)
		authorIDs = append(authorIDs, id)
	}

	// fetch all data in parallel
	var (
		authors          []sqlc.User
		categories       []sqlc.GetCategoriesByPostIDsRow
		likeCounts       []sqlc.GetLikeCountsByPostIDsRow
		commentCounts    []sqlc.GetCommentCountsByPostIDsRow
		userLikedPostIDs []pgtype.UUID
		mu               sync.Mutex
	)

	g, gCtx := errgroup.WithContext(ctx)

	// Fetch authors
	g.Go(func() error {
		result, err := repo.GetUsersByIDs(gCtx, authorIDs)
		if err != nil {
			return fmt.Errorf("failed to get authors: %w", err)
		}
		mu.Lock()
		authors = result
		mu.Unlock()
		return nil
	})

	// Fetch categories
	g.Go(func() error {
		result, err := repo.GetCategoriesByPostIDs(gCtx, postIDs)
		if err != nil {
			return fmt.Errorf("failed to get categories: %w", err)
		}
		mu.Lock()
		categories = result
		mu.Unlock()
		return nil
	})

	// Fetch like counts
	g.Go(func() error {
		result, err := repo.GetLikeCountsByPostIDs(gCtx, postIDs)
		if err != nil {
			return fmt.Errorf("failed to get like counts: %w", err)
		}
		mu.Lock()
		likeCounts = result
		mu.Unlock()
		return nil
	})

	// Fetch comment counts
	g.Go(func() error {
		result, err := repo.GetCommentCountsByPostIDs(gCtx, postIDs)
		if err != nil {
			return fmt.Errorf("failed to get comment counts: %w", err)
		}
		mu.Lock()
		commentCounts = result
		mu.Unlock()
		return nil
	})

	// Fetch user liked post IDs if authenticated
	if requestingUserID != nil && requestingUserID.Valid {
		g.Go(func() error {
			result, err := repo.GetUserLikedPostIDs(gCtx, sqlc.GetUserLikedPostIDsParams{
				UserID:  *requestingUserID,
				Column2: postIDs,
			})
			if err != nil {
				return fmt.Errorf("failed to get user liked post IDs: %w", err)
			}
			mu.Lock()
			userLikedPostIDs = result
			mu.Unlock()
			return nil
		})
	}

	// wait for all fetches to complete
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// build lookup maps
	authorMap := make(map[string]sqlc.User)
	for _, author := range authors {
		authorMap[author.ID.String()] = author
	}

	categoryMap := make(map[string][]CategoryDTO)
	for _, cat := range categories {
		postIDStr := cat.PostID.String()
		categoryMap[postIDStr] = append(categoryMap[postIDStr], CategoryDTO{
			ID:        cat.ID.String(),
			Name:      cat.Name,
			Slug:      cat.Slug,
			CreatedAt: cat.CreatedAt.Time,
		})
	}

	likeCountMap := make(map[string]int64)
	for _, likeCount := range likeCounts {
		likeCountMap[likeCount.PostID.String()] = likeCount.LikeCount
	}

	commentCountMap := make(map[string]int64)
	for _, commentCount := range commentCounts {
		commentCountMap[commentCount.PostID.String()] = commentCount.CommentCount
	}

	userLikedPostIDMap := make(map[string]bool)
	for _, likedPostID := range userLikedPostIDs {
		userLikedPostIDMap[likedPostID.String()] = true
	}

	// assemble final DTOs
	result := make([]PostCardDTO, len(posts))
	for i, post := range posts {
		postIDStr := post.ID.String()
		authorIDStr := post.AuthorID.String()

		author := authorMap[authorIDStr]

		result[i] = PostCardDTO{
			ID:          postIDStr,
			AuthorID:    authorIDStr,
			Title:       post.Title,
			Slug:        post.Slug,
			Body:        post.Body,
			IsPublished: post.IsPublished,
			CreatedAt:   post.CreatedAt.Time,
			UpdatedAt:   post.UpdatedAt.Time,
			Author: AuthorDTO{
				ID:        author.ID.String(),
				GoogleID:  author.GoogleID,
				Username:  author.Username.String,
				Email:     author.Email,
				AvatarURL: author.AvatarUrl.String,
				CreatedAt: author.CreatedAt.Time,
				UpdatedAt: author.UpdatedAt.Time,
			},
			Categories:   categoryMap[postIDStr],
			LikeCount:    likeCountMap[postIDStr],
			CommentCount: commentCountMap[postIDStr],
			UserHasLiked: userLikedPostIDMap[postIDStr],
		}

		if result[i].Categories == nil {
			result[i].Categories = []CategoryDTO{}
		}
	}

	return result, nil
}

func EnrichCommentsWithAuthors(ctx context.Context, repo *sqlc.Queries, comments []sqlc.Comment) ([]CommentDTO, error) {
	if len(comments) == 0 {
		return []CommentDTO{}, nil
	}

	// extract unique author IDs
	authorIDmap := make(map[string]bool)

	for _, comment := range comments {
		authorIDmap[comment.UserID.String()] = true
	}

	authorIDs := make([]pgtype.UUID, 0, len(authorIDmap))
	for authorID := range authorIDmap {
		id, _ := utils.StrToUUID(authorID)
		authorIDs = append(authorIDs, id)
	}

	// fetch authors
	authors, err := repo.GetUsersByIDs(ctx, authorIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get authors: %w", err)
	}

	// build author lookup map
	authorMap := make(map[string]sqlc.User)
	for _, author := range authors {
		authorMap[author.ID.String()] = author
	}

	// assemble final DTOs
	result := make([]CommentDTO, len(comments))
	for i, comment := range comments {
		authorIDStr := comment.UserID.String()
		author := authorMap[authorIDStr]
		result[i] = CommentDTO{
			ID:              comment.ID.String(),
			PostID:          comment.PostID.String(),
			UserID:          authorIDStr,
			ParentCommentID: comment.ParentCommentID.String(),
			Body:            comment.Body,
			CreatedAt:       comment.CreatedAt.Time,
			UpdatedAt:       comment.UpdatedAt.Time,
			Author: AuthorDTO{
				ID:        author.ID.String(),
				GoogleID:  author.GoogleID,
				Username:  author.Username.String,
				Email:     author.Email,
				AvatarURL: author.AvatarUrl.String,
				CreatedAt: author.CreatedAt.Time,
				UpdatedAt: author.UpdatedAt.Time,
			},
		}
	}
	return result, nil
}

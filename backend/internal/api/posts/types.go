package posts

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
)

type Service interface {
	getPostsPaginated(ctx context.Context, search string, categorySlug string, limit, offset int, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error)
	getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
	getPostBySlug(ctx context.Context, slug string, requestingUserID *pgtype.UUID) (common.PostCardDTO, error)
	getCommentsByPostSlug(ctx context.Context, slug string) ([]common.CommentDTO, error)
	togglePostLike(ctx context.Context, postID pgtype.UUID, userID pgtype.UUID) (bool, error)
	CreatePost(ctx context.Context, title string, body string, authorID string, categoryIDs []string) (common.PostCardDTO, error)
}

type PaginatedResponse struct {
	Posts   []common.PostCardDTO `json:"posts"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	HasMore bool                 `json:"hasMore"`
}

type PostLikeRequest struct {
	PostID pgtype.UUID `json:"post_id"`
}

type CreatePostRequest struct {
	Title       string   `json:"title" validate:"required,min=3,max=200"`
	Body        string   `json:"body" validate:"required,min=10"`
	CategoryIDs []string `json:"category_ids" validate:"required,min=1,max=3,dive,uuid"`
}

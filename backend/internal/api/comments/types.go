package comments

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
)

type Service interface {
	CreateComment(ctx context.Context, postID string, userID string, body string) (common.CommentDTO, error)
	DeleteComment(ctx context.Context, commentID string) error
	UpdateComment(ctx context.Context, commentID pgtype.UUID, body string) (common.CommentDTO, error)
}

type CreateCommentRequest struct {
	PostID string `json:"post_id"`
	Body   string `json:"body"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}

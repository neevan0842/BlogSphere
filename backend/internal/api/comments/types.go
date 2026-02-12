package comments

import (
	"context"

	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type Service interface {
	HandleCreateComment(ctx context.Context, postID string, userID string, body string) (sqlc.Comment, error)
	HandleDeleteComment(ctx context.Context, commentID string) error
}

type CreateCommentRequest struct {
	PostID string `json:"post_id"`
	Body   string `json:"body"`
}

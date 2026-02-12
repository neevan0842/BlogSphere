package posts

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
)

type Service interface {
	getPostsPaginated(ctx context.Context, search string, limit, offset int, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error)
	getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
	getPostBySlug(ctx context.Context, slug string, requestingUserID *pgtype.UUID) (common.PostCardDTO, error)
}

type PaginatedResponse struct {
	Posts   []common.PostCardDTO `json:"posts"`
	Page    int                  `json:"page"`
	Limit   int                  `json:"limit"`
	HasMore bool                 `json:"hasMore"`
}

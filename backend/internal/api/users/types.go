package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type Service interface {
	getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
	getUserByUsername(ctx context.Context, username pgtype.Text) (sqlc.User, error)
}

type GetUserRequest struct {
	UserID string `json:"user_id"`
}

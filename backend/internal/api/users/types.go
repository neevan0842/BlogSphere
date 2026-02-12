package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/internal/common"
)

type Service interface {
	getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
	getUserByUsername(ctx context.Context, username pgtype.Text) (sqlc.User, error)
	updateUserDescription(ctx context.Context, userID pgtype.UUID, description pgtype.Text) (sqlc.User, error)
	getPostsByUsername(ctx context.Context, username pgtype.Text, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error)
	getLikedPostsByUsername(ctx context.Context, username pgtype.Text, requestingUserID *pgtype.UUID) ([]common.PostCardDTO, error)
	deleteUserByID(ctx context.Context, userID pgtype.UUID) error
}

type UpdateUserRequest struct {
	Description string `json:"description"`
}

package common

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/utils"
)

// GetRequestingUserID extracts and validates the user ID from the request token
func GetRequestingUserID(ctx context.Context, w http.ResponseWriter, r *http.Request, repo *sqlc.Queries) *pgtype.UUID {
	token := utils.GetTokenFromRequest(r)
	userIDUUID, err := utils.GetUserIDFromToken(w, token)
	if err != nil {
		return nil
	}

	user, err := repo.GetUserByID(ctx, userIDUUID)
	if err != nil {
		return nil
	}

	return &user.ID
}

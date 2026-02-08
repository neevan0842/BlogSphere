package auth

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type Service interface {
	generateStateOauthCookie(w http.ResponseWriter) string
	getUserDataFromGoogle(code string) (sqlc.CreateOrUpdateUserParams, error)
	createOrUpdateUser(ctx context.Context, arg sqlc.CreateOrUpdateUserParams) (sqlc.CreateOrUpdateUserRow, error)
	GetUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
}

type GoogleUserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

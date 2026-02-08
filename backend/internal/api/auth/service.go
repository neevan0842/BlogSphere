package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/utils"
)

type svc struct {
	repo *sqlc.Queries
	db   *pgxpool.Pool
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func NewService(repo *sqlc.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) generateStateOauthCookie(w http.ResponseWriter) string {
	var expirationInMinutes int64 = 10
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	utils.SetCookie(w, "oauthstate", state, expirationInMinutes)

	return state
}

func (s *svc) getUserDataFromGoogle(code string) (sqlc.CreateOrUpdateUserParams, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return sqlc.CreateOrUpdateUserParams{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return sqlc.CreateOrUpdateUserParams{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return sqlc.CreateOrUpdateUserParams{}, fmt.Errorf("failed read response: %s", err.Error())
	}

	// parse google user data
	data, err := parseGoogleUserData(contents)
	if err != nil {
		return sqlc.CreateOrUpdateUserParams{}, fmt.Errorf("failed to parse user data: %s", err.Error())
	}

	// convert to CreateOrUpdateUserParams
	CreateOrUpdateUserParams := convertToCreateOrUpdateUserParams(data)

	return CreateOrUpdateUserParams, nil
}

func (s *svc) createOrUpdateUser(ctx context.Context, arg sqlc.CreateOrUpdateUserParams) (sqlc.CreateOrUpdateUserRow, error) {
	user, err := s.repo.CreateOrUpdateUser(ctx, arg)
	if err != nil {
		return sqlc.CreateOrUpdateUserRow{}, fmt.Errorf("failed to create or update user: %s", err.Error())
	}
	return user, nil
}

func parseGoogleUserData(data []byte) (*GoogleUserResponse, error) {
	var googleUser GoogleUserResponse
	err := json.Unmarshal(data, &googleUser)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %s", err.Error())
	}
	return &googleUser, nil
}

func convertToCreateOrUpdateUserParams(googleUser *GoogleUserResponse) sqlc.CreateOrUpdateUserParams {
	return sqlc.CreateOrUpdateUserParams{
		GoogleID: googleUser.ID,
		Username: pgtype.Text{
			String: googleUser.Name,
			Valid:  googleUser.Name != "",
		},
		Email: googleUser.Email,
		AvatarUrl: pgtype.Text{
			String: googleUser.Picture,
			Valid:  googleUser.Picture != "",
		},
	}
}

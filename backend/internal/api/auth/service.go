package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func (s *svc) getUserDataFromGoogle(code string) (sqlc.CreateUserParams, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return sqlc.CreateUserParams{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return sqlc.CreateUserParams{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return sqlc.CreateUserParams{}, fmt.Errorf("failed read response: %s", err.Error())
	}

	// parse google user data
	data, err := parseGoogleUserData(contents)
	if err != nil {
		return sqlc.CreateUserParams{}, fmt.Errorf("failed to parse user data: %s", err.Error())
	}

	// convert to CreateUserParams
	CreateUserParams := convertToCreateUserParams(data)

	return CreateUserParams, nil
}

func (s *svc) createUserIfNotExists(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByGoogleID(ctx, arg.GoogleID)
	if err == nil {
		return existingUser, nil
	}

	// Create new user
	user, err := s.repo.CreateUser(ctx, arg)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to create user: %s", err.Error())
	}
	return user, nil
}

func (s *svc) GetUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("failed to get user by ID: %s", err.Error())
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

func convertToCreateUserParams(googleUser *GoogleUserResponse) sqlc.CreateUserParams {
	// Extract username from email (part before @)
	username := googleUser.Email
	if atIndex := strings.Index(googleUser.Email, "@"); atIndex != -1 {
		username = googleUser.Email[:atIndex]
	}

	return sqlc.CreateUserParams{
		GoogleID: googleUser.ID,
		Username: pgtype.Text{
			String: username,
			Valid:  username != "",
		},
		Email: googleUser.Email,
		AvatarUrl: pgtype.Text{
			String: googleUser.Picture,
			Valid:  googleUser.Picture != "",
		},
	}
}

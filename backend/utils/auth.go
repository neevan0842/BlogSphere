package utils

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/config"
)

type contextKey string

const UserContextKey contextKey = "userID"

func CreateJWT(userID string, expirationInMinutes int64) (string, error) {
	secret := []byte(config.Envs.JWT_SECRET)
	expiration := time.Minute * time.Duration(expirationInMinutes)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"exp":    time.Now().Add(expiration).Unix(),
			"iat":    time.Now().Unix(),
		})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func SetCookie(w http.ResponseWriter, name string, value string, expirationInMinutes int64) {
	expiration := time.Minute * time.Duration(expirationInMinutes)
	expirationTime := time.Now().Add(expiration)

	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   config.Envs.Secure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}

	// Set the cookie in the HTTP response
	http.SetCookie(w, &cookie)
}

func GetAccessAndRefreshTokens(userID string) (string, string, error) {
	accessToken, err := CreateJWT(userID, config.Envs.ACCESS_TOKEN_EXPIRE_MINUTES)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := CreateJWT(userID, config.Envs.REFRESH_TOKEN_EXPIRE_MINUTES)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func PermissionDenied(w http.ResponseWriter) {
	WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromToken(w http.ResponseWriter, tokenString string) (pgtype.UUID, error) {
	if tokenString == "" {
		PermissionDenied(w)
		return pgtype.UUID{}, fmt.Errorf("missing token")
	}
	token, err := validateJWT(tokenString)
	if err != nil || !token.Valid {
		PermissionDenied(w)
		return pgtype.UUID{}, fmt.Errorf("invalid token: %s", err.Error())
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, ok := claims["userID"].(string)
	if !ok {
		PermissionDenied(w)
		return pgtype.UUID{}, fmt.Errorf("invalid token claims: userID not found")
	}
	userIDUUID, err := StrToUUID(userID)
	if err != nil {
		PermissionDenied(w)
		return pgtype.UUID{}, fmt.Errorf("invalid userID in token: %s", err.Error())
	}
	return userIDUUID, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWT_SECRET), nil
	})
}

func GetTokenFromRequest(r *http.Request) string {
	// Check Authorization header first (for API clients like Postman)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		// Expected format: "Bearer <token>"
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			return authHeader[7:]
		}
	}
	return ""
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserContextKey).(string)
	return userID, ok
}

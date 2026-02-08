package utils

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neevan0842/BlogSphere/backend/config"
)

func CreateJWT(userID string, expirationInMinutes int64) (string, error) {
	secret := []byte(config.Envs.JWT_SECRET)
	expiration := time.Minute * time.Duration(expirationInMinutes)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userID": userID,
			"exp":    time.Now().Add(expiration).Unix(),
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
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
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

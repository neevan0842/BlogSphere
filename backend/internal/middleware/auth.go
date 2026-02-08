package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/utils"
	"go.uber.org/zap"
)

type middleware struct {
	repo   *sqlc.Queries
	logger *zap.SugaredLogger
}

func NewMiddleware(repo *sqlc.Queries, logger *zap.SugaredLogger) *middleware {
	return &middleware{
		repo:   repo,
		logger: logger,
	}
}

func (m *middleware) UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r, "access_token")
		if tokenString == "" {
			m.logger.Error("missing token")
			utils.PermissionDenied(w)
			return
		}
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			m.logger.Errorf("invalid token: %s", err.Error())
			utils.PermissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["userID"].(string)
		if !ok {
			m.logger.Error("invalid token claims: userID not found")
			utils.PermissionDenied(w)
			return
		}
		userIDUUID, err := utils.StrToUUID(userID)
		if err != nil {
			m.logger.Errorf("invalid userID in token: %s", err.Error())
			utils.PermissionDenied(w)
			return
		}

		// get user from database and add it to context
		user, err := m.repo.GetUserByID(r.Context(), userIDUUID)
		if err != nil {
			m.logger.Errorf("failed to get user from database: %s", err.Error())
			utils.PermissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserContextKey, user.ID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"context"
	"fmt"
	"net/http"

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
		userIDUUID, err := utils.GetUserIDFromToken(w, r, "access_token")
		if err != nil {
			m.logger.Errorf("invalid userID in token: %s", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		// get user from database and add it to context
		user, err := m.repo.GetUserByID(r.Context(), userIDUUID)
		if err != nil {
			m.logger.Errorf("failed to get user from database: %s", err.Error())
			utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.UserContextKey, user.ID.String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

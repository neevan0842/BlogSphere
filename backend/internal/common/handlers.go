package common

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
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

func ExecTx(ctx context.Context, pool *pgxpool.Pool, fn func(*sqlc.Queries) error) error {
	// begin a new transaction
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// ensure we rollback if anything goes wrong
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// create a Queries instance that uses this transaction
	q := sqlc.New(tx)

	// run the transactional logic
	if err := fn(q); err != nil {
		// rollback if the function returns an error
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx %v, rollback %v", err, rbErr)
		}
		return err
	}

	// commit if no errors
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

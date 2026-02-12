package comments

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
	"github.com/neevan0842/BlogSphere/backend/utils"
)

type svc struct {
	repo *sqlc.Queries
	db   *pgxpool.Pool
}

func NewService(repo *sqlc.Queries, db *pgxpool.Pool) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) HandleCreateComment(ctx context.Context, postID string, userID string, body string) (sqlc.Comment, error) {
	postIDUUID, _ := utils.StrToUUID(postID)
	userIDUUID, _ := utils.StrToUUID(userID)

	comment, err := s.repo.CreateComment(ctx, sqlc.CreateCommentParams{
		PostID: postIDUUID,
		UserID: userIDUUID,
		Body:   body,
	})

	return comment, err
}

func (s *svc) HandleDeleteComment(ctx context.Context, commentID string) error {
	commentIDUUID, _ := utils.StrToUUID(commentID)

	return s.repo.DeleteComment(ctx, commentIDUUID)
}

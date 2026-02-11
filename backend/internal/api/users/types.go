package users

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/neevan0842/BlogSphere/backend/database/sqlc"
)

type Service interface {
	getUserByID(ctx context.Context, userID pgtype.UUID) (sqlc.User, error)
	getUserByUsername(ctx context.Context, username pgtype.Text) (sqlc.User, error)
	updateUserDescription(ctx context.Context, userID pgtype.UUID, description pgtype.Text) (sqlc.User, error)
	getPostsByUsername(ctx context.Context, username pgtype.Text) ([]UserPostDTO, error)
	getLikedPostsByUsername(ctx context.Context, username pgtype.Text, requestingUserID *pgtype.UUID) ([]LikedPostDTO, error)
}

type UpdateUserRequest struct {
	Description string `json:"description"`
}

// CategoryDTO represents category information in post responses
type CategoryDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthorDTO represents author information in post responses
type AuthorDTO struct {
	ID        string    `json:"id"`
	GoogleID  string    `json:"google_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserPostDTO represents a post by a user with enriched data
type UserPostDTO struct {
	ID           string        `json:"id"`
	AuthorID     string        `json:"author_id"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Body         string        `json:"body"`
	IsPublished  bool          `json:"is_published"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Author       AuthorDTO     `json:"author"`
	Categories   []CategoryDTO `json:"categories"`
	LikeCount    int64         `json:"like_count"`
	CommentCount int64         `json:"comment_count"`
}

// LikedPostDTO represents a post liked by a user with enriched data
type LikedPostDTO struct {
	ID           string        `json:"id"`
	AuthorID     string        `json:"author_id"`
	Title        string        `json:"title"`
	Slug         string        `json:"slug"`
	Body         string        `json:"body"`
	IsPublished  bool          `json:"is_published"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Author       AuthorDTO     `json:"author"`
	Categories   []CategoryDTO `json:"categories"`
	LikeCount    int64         `json:"like_count"`
	CommentCount int64         `json:"comment_count"`
	UserHasLiked bool          `json:"user_has_liked"`
}

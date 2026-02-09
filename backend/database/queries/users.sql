-- name: GetUserByID :one
SELECT id, google_id, username, email, avatar_url, created_at, updated_at FROM users 
WHERE id = $1;

-- name: GetUserByGoogleID :one
SELECT id, google_id, username, email, avatar_url, created_at, updated_at FROM users 
WHERE google_id = $1;

-- name: CreateUser :one
INSERT INTO users (google_id, username, email, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, google_id, username, email, avatar_url, created_at, updated_at;

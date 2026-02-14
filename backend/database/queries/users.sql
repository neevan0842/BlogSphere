-- name: GetUserByID :one
SELECT id, google_id, username, email, description, avatar_url, created_at, updated_at FROM users 
WHERE id = $1;

-- name: GetUsersByIDs :many
SELECT *
FROM users
WHERE id = ANY($1::uuid[]);

-- name: GetUserByGoogleID :one
SELECT id, google_id, username, email, description, avatar_url, created_at, updated_at FROM users 
WHERE google_id = $1;

-- name: GetUserByUsername :one
SELECT id, google_id, username, email, description, avatar_url, created_at, updated_at FROM users 
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (google_id, username, email, avatar_url)
VALUES ($1, $2, $3, $4)
RETURNING id, google_id, username, email, description, avatar_url, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET description = $2, updated_at = now()
WHERE id = $1
RETURNING id, google_id, username, email, description, avatar_url, created_at, updated_at;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;
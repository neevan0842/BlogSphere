-- name: GetUserByEmail :one
SELECT id, google_id, username, email, avatar_url, created_at, updated_at FROM users 
WHERE email = $1;

-- name: CreateOrUpdateUser :one
INSERT INTO users (google_id, username, email, avatar_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (google_id) 
DO UPDATE SET 
    username = EXCLUDED.username,
    email = EXCLUDED.email,
    avatar_url = EXCLUDED.avatar_url,
    updated_at = now()
RETURNING id, username, email, avatar_url, created_at, updated_at;
-- name: GetCommentsByPostSlug :many
SELECT c.*
FROM comments c
JOIN posts p ON p.id = c.post_id
WHERE p.slug = $1
ORDER BY c.created_at DESC;

-- name: CreateComment :one
INSERT INTO comments (post_id, user_id, body)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateComment :one
UPDATE comments
SET body = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;

-- name: GetCommentByID :one
SELECT * FROM comments WHERE id = $1;

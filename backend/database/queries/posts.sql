-- name: GetPostsByUsername :many
SELECT p.*
FROM posts p
JOIN users u ON u.id = p.author_id
WHERE u.username = $1
ORDER BY p.created_at DESC;

-- name: GetPostsLikedByUsername :many
SELECT p.*
FROM post_likes pl
JOIN users u ON u.id = pl.user_id
JOIN posts p ON p.id = pl.post_id
WHERE u.username = $1
ORDER BY p.created_at DESC;

-- name: GetUsersByIDs :many
SELECT *
FROM users
WHERE id = ANY($1::uuid[]);

-- name: GetCategoriesByPostIDs :many
SELECT
    pc.post_id,
    c.id,
    c.name,
    c.slug,
    c.created_at
FROM post_categories pc
JOIN categories c ON c.id = pc.category_id
WHERE pc.post_id = ANY($1::uuid[]);

-- name: GetLikeCountsByPostIDs :many
SELECT
    post_id,
    COUNT(*)::bigint AS like_count
FROM post_likes
WHERE post_id = ANY($1::uuid[])
GROUP BY post_id;

-- name: GetCommentCountsByPostIDs :many
SELECT
    post_id,
    COUNT(*)::bigint AS comment_count
FROM comments
WHERE post_id = ANY($1::uuid[])
GROUP BY post_id;


-- name: GetUserLikedPostIDs :many
SELECT post_id
FROM post_likes
WHERE user_id = $1
AND post_id = ANY($2::uuid[]);

-- name: GetPostBySearchPaginated :many
SELECT p.*
FROM posts p
WHERE p.title ILIKE '%' || COALESCE($1, '') || '%'
ORDER BY p.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPostBySlug :one
SELECT * 
FROM posts 
WHERE slug = $1;

-- name: CreatePostLike :one
INSERT INTO post_likes (post_id, user_id)
VALUES ($1, $2)
ON CONFLICT (post_id, user_id) DO NOTHING
RETURNING *;

-- name: DeletePostLike :exec
DELETE FROM post_likes
WHERE post_id = $1 AND user_id = $2;

-- name: GetPostLike :one
SELECT *
FROM post_likes
WHERE post_id = $1 AND user_id = $2;
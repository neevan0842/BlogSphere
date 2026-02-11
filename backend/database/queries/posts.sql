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
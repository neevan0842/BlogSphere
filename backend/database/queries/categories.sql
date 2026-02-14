-- name: GetCategories :many
SELECT * FROM categories ORDER BY name;

-- name: BatchCreatePostCategories :exec
INSERT INTO post_categories (post_id, category_id)
SELECT unnest(sqlc.narg('post_id')::uuid[]), unnest(sqlc.narg('category_id')::uuid[]);

-- name: DeletePostCategoriesByPostID :exec
DELETE FROM post_categories WHERE post_id = $1;
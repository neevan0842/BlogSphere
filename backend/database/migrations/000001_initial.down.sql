-- Drop indexes
DROP INDEX IF EXISTS idx_comment_post_id;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS user_follows;
DROP TABLE IF EXISTS comment_likes;
DROP TABLE IF EXISTS post_likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS post_categories;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;

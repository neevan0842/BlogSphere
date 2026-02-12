CREATE INDEX idx_post_likes_post_id ON post_likes(post_id);
CREATE INDEX idx_post_likes_user_id ON post_likes(user_id);
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_post_categories_post_id ON post_categories(post_id);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
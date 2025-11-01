-- Create post_versions table for the hybrid versioning system
CREATE TABLE IF NOT EXISTS post_versions (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content_json TEXT NOT NULL,
    content_markdown TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('draft', 'published')) DEFAULT 'draft',
    author_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add new fields to posts table
ALTER TABLE posts ADD COLUMN IF NOT EXISTS slug VARCHAR(255);
ALTER TABLE posts ADD COLUMN IF NOT EXISTS published_version_id INTEGER REFERENCES post_versions(id);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_post_versions_post_id ON post_versions(post_id);
CREATE INDEX IF NOT EXISTS idx_post_versions_status ON post_versions(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_version_id ON posts(published_version_id);

-- Migrate existing post data to versions
-- For each existing post, create an initial version with the current content and status
INSERT INTO post_versions (post_id, title, content_json, content_markdown, status, author_id)
SELECT
    p.id,
    p.title,
    COALESCE(p.content_json, '{}'),
    COALESCE(p.content_markdown, ''),
    CASE
        WHEN p.status IS NULL OR p.status = '' THEN 'draft'
        ELSE LOWER(p.status::text)
    END,
    p.user_id
FROM posts p;

-- Update posts to set published_version_id for posts that had 'published' status
UPDATE posts
SET published_version_id = pv.id
FROM post_versions pv
WHERE posts.id = pv.post_id AND pv.status = 'published';

-- Remove old status and content fields from posts table
ALTER TABLE posts DROP COLUMN IF EXISTS status;
ALTER TABLE posts DROP COLUMN IF EXISTS content_json;
ALTER TABLE posts DROP COLUMN IF EXISTS content_markdown;

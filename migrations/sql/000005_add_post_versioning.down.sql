-- Reverse migration for post versioning system

-- Add back the removed columns to posts table
ALTER TABLE posts ADD COLUMN status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published'));
ALTER TABLE posts ADD COLUMN content_json TEXT;
ALTER TABLE posts ADD COLUMN content_markdown TEXT;

-- Migrate version data back to posts table
UPDATE posts
SET
    status = pv.status,
    content_json = pv.content_json,
    content_markdown = pv.content_markdown,
    title = pv.title
FROM post_versions pv
WHERE posts.id = pv.post_id
    AND (posts.published_version_id = pv.id OR (posts.published_version_id IS NULL AND pv.status = 'draft'));

-- Drop indexes
DROP INDEX IF EXISTS idx_posts_published_version_id;
DROP INDEX IF EXISTS idx_post_versions_status;
DROP INDEX IF EXISTS idx_post_versions_post_id;

-- Remove added fields from posts table
ALTER TABLE posts DROP COLUMN IF EXISTS slug;
ALTER TABLE posts DROP COLUMN IF EXISTS published_version_id;

-- Drop post_versions table
DROP TABLE IF EXISTS post_versions;

-- Recreate any dropped indexes if they don't exist
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

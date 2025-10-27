-- Reverse migration for post versioning system

-- Add back the removed columns to posts table
ALTER TABLE posts ADD COLUMN status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published'));
ALTER TABLE posts ADD COLUMN content_json TEXT;
ALTER TABLE posts ADD COLUMN content_markdown TEXT;

-- Migrate version data back to posts table
-- Use the published version if exists, otherwise use the first draft version
UPDATE posts p
SET
    status = COALESCE(
        (SELECT pv.status FROM post_versions pv WHERE pv.post_id = p.id AND pv.status = 'published' ORDER BY pv.created_at DESC LIMIT 1),
        (SELECT pv.status FROM post_versions pv WHERE pv.post_id = p.id ORDER BY pv.created_at ASC LIMIT 1)
    ),
    content_json = COALESCE(
        (SELECT pv.content_json FROM post_versions pv WHERE pv.post_id = p.id AND pv.status = 'published' ORDER BY pv.created_at DESC LIMIT 1),
        (SELECT pv.content_json FROM post_versions pv WHERE pv.post_id = p.id ORDER BY pv.created_at ASC LIMIT 1)
    ),
    content_markdown = COALESCE(
        (SELECT pv.content_markdown FROM post_versions pv WHERE pv.post_id = p.id AND pv.status = 'published' ORDER BY pv.created_at DESC LIMIT 1),
        (SELECT pv.content_markdown FROM post_versions pv WHERE pv.post_id = p.id ORDER BY pv.created_at ASC LIMIT 1)
    ),
    title = COALESCE(
        (SELECT pv.title FROM post_versions pv WHERE pv.post_id = p.id AND pv.status = 'published' ORDER BY pv.created_at DESC LIMIT 1),
        (SELECT pv.title FROM post_versions pv WHERE pv.post_id = p.id ORDER BY pv.created_at ASC LIMIT 1)
    );

-- Drop indexes
DROP INDEX IF EXISTS idx_post_versions_status;
DROP INDEX IF EXISTS idx_post_versions_post_id;

-- Remove added fields from posts table
ALTER TABLE posts DROP COLUMN IF EXISTS slug;

-- Drop post_versions table
DROP TABLE IF EXISTS post_versions;

-- Recreate any dropped indexes if they don't exist
CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at DESC);

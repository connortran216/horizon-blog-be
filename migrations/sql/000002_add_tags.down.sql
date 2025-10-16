-- Drop indexes
DROP INDEX IF EXISTS idx_post_tags_tag_id;
DROP INDEX IF EXISTS idx_post_tags_post_id;
DROP INDEX IF EXISTS idx_tags_usage_count;
DROP INDEX IF EXISTS idx_tags_name;

-- Drop tables (order matters due to foreign key constraints)
DROP TABLE IF EXISTS post_tags;
DROP TABLE IF EXISTS tags;

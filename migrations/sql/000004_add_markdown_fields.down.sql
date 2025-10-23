-- Rollback: restore old content column and remove new ones
ALTER TABLE posts ADD COLUMN content TEXT;
ALTER TABLE posts DROP COLUMN IF EXISTS content_markdown;
ALTER TABLE posts DROP COLUMN IF EXISTS content_json;

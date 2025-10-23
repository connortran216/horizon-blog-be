-- Add markdown and ProseMirror JSON content fields
ALTER TABLE posts ADD COLUMN content_markdown TEXT;
ALTER TABLE posts ADD COLUMN content_json TEXT;

-- Drop old content column (this will clear existing posts as requested)
ALTER TABLE posts DROP COLUMN IF EXISTS content;

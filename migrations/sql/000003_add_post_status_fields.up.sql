-- Add is_draft and is_published columns to posts table
ALTER TABLE posts ADD COLUMN is_draft BOOLEAN DEFAULT true NOT NULL;
ALTER TABLE posts ADD COLUMN is_published BOOLEAN DEFAULT false NOT NULL;

-- Update existing posts to be published since they're already created
UPDATE posts SET is_draft = false, is_published = true;

-- Add check constraint to ensure a post cannot be both draft and published (optional, but recommended)
-- Uncomment if you want to enforce mutual exclusivity:
-- ALTER TABLE posts ADD CONSTRAINT check_post_status
-- CHECK (NOT (is_draft = true AND is_published = true));

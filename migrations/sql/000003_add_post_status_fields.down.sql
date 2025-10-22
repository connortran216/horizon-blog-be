-- Drop the check constraint if it was added (uncomment if used)
-- ALTER TABLE posts DROP CONSTRAINT IF EXISTS check_post_status;

-- Drop the columns
ALTER TABLE posts DROP COLUMN IF EXISTS is_published;
ALTER TABLE posts DROP COLUMN IF EXISTS is_draft;

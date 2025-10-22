-- Drop the status column
ALTER TABLE posts DROP COLUMN IF EXISTS status;

-- Drop the enum type
DROP TYPE IF EXISTS post_status;

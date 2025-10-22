-- Create enum type for post status
CREATE TYPE post_status AS ENUM ('draft', 'published');

-- Add status column to posts table
ALTER TABLE posts ADD COLUMN status post_status DEFAULT 'draft' NOT NULL;

-- Migrate existing posts to appropriate status
-- Since this replaces the previous migration, we check for existing boolean columns
UPDATE posts SET status = 'published' WHERE is_published = true;
UPDATE posts SET status = 'draft' WHERE is_published = false OR is_published IS NULL;


-- Initial database setup for Go CRUD API
-- This file runs when the PostgreSQL container starts for the first time

-- Create the database if it doesn't exist
-- (This is handled by POSTGRES_DB environment variable)

-- Enable UUID extension if needed (optional, for future features)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create indexes for better performance (these will be created by GORM auto-migration)
-- But we can add specific optimized indexes here if needed

-- Set timezone
SET timezone = 'UTC';

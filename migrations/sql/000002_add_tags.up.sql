-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create post_tags junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, tag_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
CREATE INDEX IF NOT EXISTS idx_tags_usage_count ON tags(usage_count DESC);
CREATE INDEX IF NOT EXISTS idx_post_tags_post_id ON post_tags(post_id);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag_id ON post_tags(tag_id);

-- Insert some popular tags
INSERT INTO tags (name, description, usage_count) VALUES
    ('golang', 'Go programming language', 0),
    ('web-development', 'Web development topics', 0),
    ('tutorial', 'Tutorial and guide content', 0),
    ('programming', 'Programming concepts and techniques', 0),
    ('database', 'Database related content', 0),
    ('api', 'API development and integration', 0),
    ('docker', 'Containerization with Docker', 0),
    ('kubernetes', 'Kubernetes orchestration', 0),
    ('microservices', 'Microservices architecture', 0),
    ('testing', 'Software testing practices', 0)
ON CONFLICT (name) DO NOTHING;

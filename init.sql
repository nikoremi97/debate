-- Database initialization script for the debate chatbot

-- Users table (for future user management)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) PRIMARY KEY, -- ULID format
    email VARCHAR(255) UNIQUE,
    name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Topics table (predefined or user-created)
CREATE TABLE IF NOT EXISTS topics (
    id VARCHAR(26) PRIMARY KEY, -- ULID format
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- e.g., 'technology', 'politics', 'sports'
    created_at TIMESTAMP DEFAULT NOW()
);

-- Conversations table
CREATE TABLE IF NOT EXISTS conversations (
    id VARCHAR(26) PRIMARY KEY, -- ULID format
    user_id VARCHAR(26) REFERENCES users(id),
    topic_id VARCHAR(26) REFERENCES topics(id),
    topic_name VARCHAR(255), -- denormalized for performance
    bot_stance VARCHAR(10) NOT NULL, -- 'PRO' or 'CON'
    title VARCHAR(255), -- auto-generated or user-defined
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    message_count INTEGER DEFAULT 0
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(26) PRIMARY KEY, -- ULID format
    conversation_id VARCHAR(26) REFERENCES conversations(id) ON DELETE CASCADE,
    role VARCHAR(10) NOT NULL, -- 'user' or 'bot'
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_topic_id ON conversations(topic_id);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON conversations(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_conversations_updated_at ON conversations(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Insert some popular debate topics with ULID
INSERT INTO topics (id, name, description, category) VALUES
('01HZ0000000000000000000001', 'Artificial Intelligence Regulation', 'Should AI be heavily regulated by governments?', 'technology'),
('01HZ0000000000000000000002', 'Remote Work vs Office Work', 'Which is more productive and beneficial?', 'business'),
('01HZ0000000000000000000003', 'Social Media Impact', 'Is social media good or bad for society?', 'technology'),
('01HZ0000000000000000000004', 'Climate Change Action', 'Should governments take more aggressive action on climate change?', 'politics'),
('01HZ0000000000000000000005', 'Universal Basic Income', 'Should governments provide UBI to all citizens?', 'politics'),
('01HZ0000000000000000000006', 'Cryptocurrency Future', 'Will cryptocurrency replace traditional money?', 'finance'),
('01HZ0000000000000000000007', 'Space Exploration', 'Should we invest more in space exploration?', 'science'),
('01HZ0000000000000000000008', 'Electric Vehicles', 'Are electric vehicles the future of transportation?', 'technology'),
('01HZ0000000000000000000009', 'Online Education', 'Is online education as effective as traditional education?', 'education'),
('01HZ000000000000000000000A', 'Privacy vs Security', 'Should privacy be sacrificed for national security?', 'politics')
ON CONFLICT (name) DO NOTHING;

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at
CREATE TRIGGER update_conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

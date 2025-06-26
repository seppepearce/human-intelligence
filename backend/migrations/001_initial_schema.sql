-- Initial schema for Human Intelligence platform
-- Migration 001: Core tables and relationships

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    avatar VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT true
);

-- Nodes table (atomic learning units)
CREATE TABLE nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    content TEXT,
    node_type VARCHAR(50) NOT NULL DEFAULT 'text',
    metadata JSONB DEFAULT '{}',
    tags TEXT[] DEFAULT '{}',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_public BOOLEAN DEFAULT true,
    vote_score INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0
);

-- Learning paths table
CREATE TABLE learning_paths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    access_level VARCHAR(50) NOT NULL DEFAULT 'public',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_path_id UUID REFERENCES learning_paths(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    vote_score INTEGER DEFAULT 0,
    fork_count INTEGER DEFAULT 0,
    completion_count INTEGER DEFAULT 0,
    tags TEXT[] DEFAULT '{}'
);

-- Junction table for nodes within learning paths
CREATE TABLE path_nodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    node_id UUID NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    position INTEGER NOT NULL,
    is_required BOOLEAN DEFAULT true,
    prerequisites UUID[] DEFAULT '{}',
    UNIQUE(path_id, position),
    UNIQUE(path_id, node_id)
);

-- TLDRs table (completion summaries)
CREATE TABLE tldrs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content VARCHAR(140) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    vote_score INTEGER DEFAULT 0,
    UNIQUE(path_id, user_id)
);

-- Votes table (for nodes, paths, and tldrs)
CREATE TABLE votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id UUID NOT NULL,
    target_type VARCHAR(20) NOT NULL CHECK (target_type IN ('node', 'path', 'tldr')),
    vote_type INTEGER NOT NULL CHECK (vote_type IN (-1, 1)),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, target_id, target_type)
);

-- User progress tracking
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    completed_nodes UUID[] DEFAULT '{}',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, path_id)
);

-- Access control for private/group paths
CREATE TABLE access_control (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    path_id UUID NOT NULL REFERENCES learning_paths(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(20) NOT NULL DEFAULT 'viewer' CHECK (role IN ('viewer', 'contributor', 'admin')),
    granted_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(path_id, user_id)
);

-- Plugin system for extensible node types
CREATE TABLE plugins (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,
    author VARCHAR(255),
    config_schema JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    is_premium BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Activity events for real-time forum view
CREATE TABLE activity_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(50) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_id UUID NOT NULL,
    path_id UUID REFERENCES learning_paths(id) ON DELETE CASCADE,
    node_id UUID REFERENCES nodes(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

CREATE INDEX idx_nodes_created_by ON nodes(created_by);
CREATE INDEX idx_nodes_created_at ON nodes(created_at);
CREATE INDEX idx_nodes_node_type ON nodes(node_type);
CREATE INDEX idx_nodes_is_public ON nodes(is_public);
CREATE INDEX idx_nodes_vote_score ON nodes(vote_score);
CREATE INDEX idx_nodes_tags ON nodes USING GIN(tags);
CREATE INDEX idx_nodes_metadata ON nodes USING GIN(metadata);

CREATE INDEX idx_learning_paths_created_by ON learning_paths(created_by);
CREATE INDEX idx_learning_paths_created_at ON learning_paths(created_at);
CREATE INDEX idx_learning_paths_access_level ON learning_paths(access_level);
CREATE INDEX idx_learning_paths_parent_path_id ON learning_paths(parent_path_id);
CREATE INDEX idx_learning_paths_vote_score ON learning_paths(vote_score);
CREATE INDEX idx_learning_paths_tags ON learning_paths USING GIN(tags);

CREATE INDEX idx_path_nodes_path_id ON path_nodes(path_id);
CREATE INDEX idx_path_nodes_node_id ON path_nodes(node_id);
CREATE INDEX idx_path_nodes_position ON path_nodes(path_id, position);

CREATE INDEX idx_tldrs_path_id ON tldrs(path_id);
CREATE INDEX idx_tldrs_user_id ON tldrs(user_id);
CREATE INDEX idx_tldrs_created_at ON tldrs(created_at);
CREATE INDEX idx_tldrs_vote_score ON tldrs(vote_score);

CREATE INDEX idx_votes_user_id ON votes(user_id);
CREATE INDEX idx_votes_target ON votes(target_id, target_type);
CREATE INDEX idx_votes_created_at ON votes(created_at);

CREATE INDEX idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX idx_user_progress_path_id ON user_progress(path_id);
CREATE INDEX idx_user_progress_completed_at ON user_progress(completed_at);

CREATE INDEX idx_access_control_path_id ON access_control(path_id);
CREATE INDEX idx_access_control_user_id ON access_control(user_id);

CREATE INDEX idx_activity_events_created_at ON activity_events(created_at);
CREATE INDEX idx_activity_events_user_id ON activity_events(user_id);
CREATE INDEX idx_activity_events_event_type ON activity_events(event_type);
CREATE INDEX idx_activity_events_path_id ON activity_events(path_id);

-- Functions to update vote scores automatically
CREATE OR REPLACE FUNCTION update_vote_score()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Update the target's vote score
        IF NEW.target_type = 'node' THEN
            UPDATE nodes SET vote_score = vote_score + NEW.vote_type WHERE id = NEW.target_id;
        ELSIF NEW.target_type = 'path' THEN
            UPDATE learning_paths SET vote_score = vote_score + NEW.vote_type WHERE id = NEW.target_id;
        ELSIF NEW.target_type = 'tldr' THEN
            UPDATE tldrs SET vote_score = vote_score + NEW.vote_type WHERE id = NEW.target_id;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        -- Handle vote changes
        IF OLD.vote_type != NEW.vote_type THEN
            IF NEW.target_type = 'node' THEN
                UPDATE nodes SET vote_score = vote_score - OLD.vote_type + NEW.vote_type WHERE id = NEW.target_id;
            ELSIF NEW.target_type = 'path' THEN
                UPDATE learning_paths SET vote_score = vote_score - OLD.vote_type + NEW.vote_type WHERE id = NEW.target_id;
            ELSIF NEW.target_type = 'tldr' THEN
                UPDATE tldrs SET vote_score = vote_score - OLD.vote_type + NEW.vote_type WHERE id = NEW.target_id;
            END IF;
        END IF;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Remove vote from score
        IF OLD.target_type = 'node' THEN
            UPDATE nodes SET vote_score = vote_score - OLD.vote_type WHERE id = OLD.target_id;
        ELSIF OLD.target_type = 'path' THEN
            UPDATE learning_paths SET vote_score = vote_score - OLD.vote_type WHERE id = OLD.target_id;
        ELSIF OLD.target_type = 'tldr' THEN
            UPDATE tldrs SET vote_score = vote_score - OLD.vote_type WHERE id = OLD.target_id;
        END IF;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger for automatic vote score updates
CREATE TRIGGER trigger_update_vote_score
    AFTER INSERT OR UPDATE OR DELETE ON votes
    FOR EACH ROW EXECUTE FUNCTION update_vote_score();

-- Function to update fork count when a path is forked
CREATE OR REPLACE FUNCTION update_fork_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.parent_path_id IS NOT NULL THEN
        UPDATE learning_paths SET fork_count = fork_count + 1 WHERE id = NEW.parent_path_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for automatic fork count updates
CREATE TRIGGER trigger_update_fork_count
    AFTER INSERT ON learning_paths
    FOR EACH ROW EXECUTE FUNCTION update_fork_count();

-- Function to update completion count when a path is completed
CREATE OR REPLACE FUNCTION update_completion_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'UPDATE' AND OLD.completed_at IS NULL AND NEW.completed_at IS NOT NULL THEN
        UPDATE learning_paths SET completion_count = completion_count + 1 WHERE id = NEW.path_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for automatic completion count updates
CREATE TRIGGER trigger_update_completion_count
    AFTER UPDATE ON user_progress
    FOR EACH ROW EXECUTE FUNCTION update_completion_count();

-- Insert default plugins
INSERT INTO plugins (id, name, description, version, author, is_premium) VALUES
('text', 'Text Content', 'Basic text-based learning content', '1.0.0', 'HI Team', false),
('video', 'Video Content', 'YouTube and video embeds', '1.0.0', 'HI Team', false),
('link', 'External Links', 'Links to external resources', '1.0.0', 'HI Team', false),
('quiz', 'Interactive Quiz', 'Multiple choice quizzes', '1.0.0', 'HI Team', true),
('code', 'Code Snippets', 'Code examples and repositories', '1.0.0', 'HI Team', true),
('latex', 'LaTeX Math', 'Mathematical expressions and formulas', '1.0.0', 'HI Team', true),
('jupyter', 'Jupyter Notebooks', 'Interactive Python notebooks', '1.0.0', 'HI Team', true),
('wolfram', 'Wolfram Integration', 'Physics and mathematical computations', '1.0.0', 'HI Team', true);

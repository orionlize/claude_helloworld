-- Supabase SQL Setup Script
-- Run this in your Supabase SQL Editor to create the necessary tables

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Collections table
CREATE TABLE IF NOT EXISTS collections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id UUID REFERENCES collections(id) ON DELETE CASCADE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Endpoints table
CREATE TABLE IF NOT EXISTS endpoints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collection_id UUID NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    url TEXT NOT NULL,
    headers JSONB DEFAULT '{}',
    body TEXT,
    description TEXT,
    request_params JSONB DEFAULT '[]',
    request_body JSONB,
    response_params JSONB DEFAULT '[]',
    response_body JSONB,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Environments table
CREATE TABLE IF NOT EXISTS environments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    variables JSONB DEFAULT '{}',
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW()),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT TIMEZONE('utc', NOW())
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);
CREATE INDEX IF NOT EXISTS idx_collections_project_id ON collections(project_id);
CREATE INDEX IF NOT EXISTS idx_collections_parent_id ON collections(parent_id);
CREATE INDEX IF NOT EXISTS idx_endpoints_collection_id ON endpoints(collection_id);
CREATE INDEX IF NOT EXISTS idx_environments_project_id ON environments(project_id);

-- Enable Row Level Security (RLS)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;
ALTER TABLE projects ENABLE ROW LEVEL SECURITY;
ALTER TABLE collections ENABLE ROW LEVEL SECURITY;
ALTER TABLE endpoints ENABLE ROW LEVEL SECURITY;
ALTER TABLE environments ENABLE ROW LEVEL SECURITY;

-- RLS Policies for users
CREATE POLICY "Users can view their own data" ON users
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can insert their own data" ON users
    FOR INSERT WITH CHECK (auth.uid() = id);

CREATE POLICY "Users can update their own data" ON users
    FOR UPDATE USING (auth.uid() = id);

-- RLS Policies for projects
CREATE POLICY "Users can view their own projects" ON projects
    FOR SELECT USING (auth.uid() = user_id);

CREATE POLICY "Users can create their own projects" ON projects
    FOR INSERT WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update their own projects" ON projects
    FOR UPDATE USING (auth.uid() = user_id);

CREATE POLICY "Users can delete their own projects" ON projects
    FOR DELETE USING (auth.uid() = user_id);

-- RLS Policies for collections (via projects)
CREATE POLICY "Users can view collections of their projects" ON collections
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = collections.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can create collections in their projects" ON collections
    FOR INSERT WITH CHECK (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = collections.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can update collections of their projects" ON collections
    FOR UPDATE USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = collections.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can delete collections of their projects" ON collections
    FOR DELETE USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = collections.project_id AND projects.user_id = auth.uid()
        )
    );

-- RLS Policies for endpoints (via collections -> projects)
CREATE POLICY "Users can view endpoints of their projects" ON endpoints
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM collections
            JOIN projects ON projects.id = collections.project_id
            WHERE collections.id = endpoints.collection_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can create endpoints in their projects" ON endpoints
    FOR INSERT WITH CHECK (
        EXISTS (
            SELECT 1 FROM collections
            JOIN projects ON projects.id = collections.project_id
            WHERE collections.id = endpoints.collection_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can update endpoints of their projects" ON endpoints
    FOR UPDATE USING (
        EXISTS (
            SELECT 1 FROM collections
            JOIN projects ON projects.id = collections.project_id
            WHERE collections.id = endpoints.collection_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can delete endpoints of their projects" ON endpoints
    FOR DELETE USING (
        EXISTS (
            SELECT 1 FROM collections
            JOIN projects ON projects.id = collections.project_id
            WHERE collections.id = endpoints.collection_id AND projects.user_id = auth.uid()
        )
    );

-- RLS Policies for environments
CREATE POLICY "Users can view environments of their projects" ON environments
    FOR SELECT USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = environments.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can create environments in their projects" ON environments
    FOR INSERT WITH CHECK (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = environments.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can update environments of their projects" ON environments
    FOR UPDATE USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = environments.project_id AND projects.user_id = auth.uid()
        )
    );

CREATE POLICY "Users can delete environments of their projects" ON environments
    FOR DELETE USING (
        EXISTS (
            SELECT 1 FROM projects WHERE projects.id = environments.project_id AND projects.user_id = auth.uid()
        )
    );

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = TIMEZONE('utc', NOW());
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Triggers to automatically update updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_collections_updated_at BEFORE UPDATE ON collections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_endpoints_updated_at BEFORE UPDATE ON endpoints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_environments_updated_at BEFORE UPDATE ON environments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

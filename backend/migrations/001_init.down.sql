-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_projects_updated_at ON projects;
DROP TRIGGER IF EXISTS update_collections_updated_at ON collections;
DROP TRIGGER IF EXISTS update_endpoints_updated_at ON endpoints;
DROP TRIGGER IF EXISTS update_environments_updated_at ON environments;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS endpoints;
DROP TABLE IF EXISTS collections;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS organizations;

-- Create organizations table
CREATE TABLE organizations
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    slug text NOT NULL,
    settings JSONB DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create users table
CREATE TABLE users
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email text NOT NULL,
    password_hash text NOT NULL,
    role text NOT NULL CHECK (role IN ('superadmin', 'admin', 'user', 'viewer')),
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    first_name text,
    last_name text,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create api_keys table
CREATE TABLE api_keys
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    hashed_key text NOT NULL,
    label text NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    revoked_at timestamp with time zone NULL
);

-- Create chats table
CREATE TABLE chats
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id BIGINT NULL REFERENCES users(id) ON DELETE SET NULL,
    title text,
    tags JSONB DEFAULT '[]'::jsonb,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create messages table
CREATE TABLE messages
(
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    role text NOT NULL CHECK (role IN ('user', 'assistant', 'system')),
    content text NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    latency INTEGER DEFAULT 0,
    token_count INTEGER DEFAULT 0,
    created_at timestamp with time zone NOT NULL DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_api_keys_org_id ON api_keys(organization_id);
CREATE INDEX idx_users_org_id ON users(organization_id);
CREATE INDEX idx_chats_org_id ON chats(organization_id);
CREATE INDEX idx_chats_user_id ON chats(user_id);
CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

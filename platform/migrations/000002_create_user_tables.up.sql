-- Add UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set timezone
-- For more information, please visit:
-- https://en.wikipedia.org/wiki/List_of_tz_database_time_zones
SET TIMEZONE="America/Toronto";

-- Create users table
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4 () PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW (),
    updated_at TIMESTAMP DEFAULT NOW (),
    username VARCHAR (255) UNIQUE NOT NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    flow_address VARCHAR(255) UNIQUE NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    login_obj JSONB NULL
);

CREATE INDEX active_id ON users (id);
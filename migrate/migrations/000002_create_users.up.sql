CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  email citext UNIQUE NOT NULL, -- case-insensitive email
  username VARCHAR(255) UNIQUE NOT NULL,
  password BYTEA NOT NULL,
  created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);

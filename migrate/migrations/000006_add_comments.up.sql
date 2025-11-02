CREATE TABLE IF NOT EXISTS comments (
  id bigserial PRIMARY KEY,
  post_id bigserial NOT NULL,
  user_id bigserial NOT NULL,
  content text NOT NULL,
  created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);
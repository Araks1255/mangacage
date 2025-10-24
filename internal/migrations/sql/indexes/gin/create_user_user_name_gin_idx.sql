CREATE INDEX IF NOT EXISTS user_user_name_gin_idx
ON users USING GIN (user_name gin_trgm_ops)
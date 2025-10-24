CREATE INDEX IF NOT EXISTS user_moderator_idx
ON users (id, moderator_id)
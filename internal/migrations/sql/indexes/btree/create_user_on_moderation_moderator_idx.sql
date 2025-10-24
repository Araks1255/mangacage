CREATE INDEX IF NOT EXISTS user_on_moderation_moderator_idx
ON users_on_moderation (id, moderator_id)
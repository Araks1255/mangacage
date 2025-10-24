CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_on_moderation_user_name
ON users_on_moderation (lower(user_name))
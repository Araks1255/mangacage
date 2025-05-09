CREATE UNIQUE INDEX IF NOT EXISTS uniq_user_name
ON users (lower(user_name))
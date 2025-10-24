CREATE UNIQUE INDEX IF NOT EXISTS uniq_title_rate
ON title_rates (title_id, user_id)
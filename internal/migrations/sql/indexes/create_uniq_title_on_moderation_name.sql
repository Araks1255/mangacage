CREATE UNIQUE INDEX IF NOT EXISTS uniq_title_on_moderation_name
ON titles_on_moderation (lower(name))
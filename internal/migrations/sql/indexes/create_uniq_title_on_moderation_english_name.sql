CREATE UNIQUE INDEX IF NOT EXISTS uniq_titles_on_moderation_english_name
ON titles_on_moderation (lower(english_name))
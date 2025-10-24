CREATE UNIQUE INDEX IF NOT EXISTS uniq_author_on_moderation_english_name
ON authors_on_moderation (lower(english_name))
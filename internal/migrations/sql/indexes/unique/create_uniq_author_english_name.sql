CREATE UNIQUE INDEX IF NOT EXISTS uniq_author_english_name
ON authors (lower(english_name))
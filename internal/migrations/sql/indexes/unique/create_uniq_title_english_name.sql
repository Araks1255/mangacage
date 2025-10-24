CREATE UNIQUE INDEX IF NOT EXISTS uniq_titles_english_name
ON titles (lower(english_name))
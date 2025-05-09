CREATE UNIQUE INDEX IF NOT EXISTS uniq_title_name
ON titles (lower(name))
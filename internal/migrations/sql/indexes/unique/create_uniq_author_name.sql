CREATE UNIQUE INDEX IF NOT EXISTS uniq_author_name
ON authors (lower(name))
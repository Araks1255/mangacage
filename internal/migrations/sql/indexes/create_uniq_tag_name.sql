CREATE UNIQUE INDEX IF NOT EXISTS uniq_tag_name
ON tags (lower(name))
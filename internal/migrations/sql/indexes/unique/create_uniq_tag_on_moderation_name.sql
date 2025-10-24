CREATE UNIQUE INDEX IF NOT EXISTS uniq_tag_on_moderation_name
ON tags_on_moderation (lower(name))
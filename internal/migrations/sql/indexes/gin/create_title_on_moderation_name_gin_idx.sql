CREATE INDEX IF NOT EXISTS title_on_moderation_name_gin_idx
ON titles_on_moderation USING GIN (name gin_trgm_ops)
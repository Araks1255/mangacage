CREATE INDEX IF NOT EXISTS author_on_moderation_name_gin_idx
ON authors_on_moderation USING GIN (name gin_trgm_ops)
CREATE INDEX IF NOT EXISTS author_original_name_gin_idx
ON authors USING GIN (original_name gin_trgm_ops)
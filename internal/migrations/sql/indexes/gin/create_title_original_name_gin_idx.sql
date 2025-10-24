CREATE INDEX IF NOT EXISTS title_original_name_gin_idx
ON titles USING GIN (original_name gin_trgm_ops)
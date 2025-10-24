CREATE INDEX IF NOT EXISTS title_name_gin_idx
ON titles USING GIN (name gin_trgm_ops)
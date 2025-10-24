CREATE INDEX IF NOT EXISTS author_name_gin_idx
ON authors USING GIN (name gin_trgm_ops)
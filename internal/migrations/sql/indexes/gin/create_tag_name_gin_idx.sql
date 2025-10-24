CREATE INDEX IF NOT EXISTS tag_name_gin_idx
ON tags USING GIN (name gin_trgm_ops)
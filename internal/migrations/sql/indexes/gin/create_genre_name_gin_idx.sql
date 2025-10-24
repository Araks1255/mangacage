CREATE INDEX IF NOT EXISTS genre_name_gin_idx
ON genres USING GIN (name gin_trgm_ops)
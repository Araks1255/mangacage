CREATE INDEX IF NOT EXISTS title_english_name_gin_idx
ON titles USING GIN (english_name gin_trgm_ops)
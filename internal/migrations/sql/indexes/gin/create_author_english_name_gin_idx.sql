CREATE INDEX IF NOT EXISTS author_english_name_gin_idx
ON authors USING GIN (english_name gin_trgm_ops)
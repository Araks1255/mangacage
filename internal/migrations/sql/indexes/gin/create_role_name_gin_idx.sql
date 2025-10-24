CREATE INDEX IF NOT EXISTS role_name_gin_idx
ON roles USING GIN (name gin_trgm_ops)
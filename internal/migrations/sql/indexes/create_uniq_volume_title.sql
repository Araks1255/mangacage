CREATE UNIQUE INDEX IF NOT EXISTS uniq_volume_title
ON volumes (lower(name), title_id)
CREATE UNIQUE INDEX IF NOT EXISTS uniq_chapter_volume
ON chapters (lower(name), volume_id)
CREATE UNIQUE INDEX IF NOT EXISTS uniq_chapter_on_moderation_volume_on_moderation
ON chapters_on_moderation (lower(name), volume_on_moderation_id)
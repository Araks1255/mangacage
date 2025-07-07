CREATE UNIQUE INDEX IF NOT EXISTS uniq_volume_on_moderation_title
ON volumes_on_moderation (lower(name), title_id)
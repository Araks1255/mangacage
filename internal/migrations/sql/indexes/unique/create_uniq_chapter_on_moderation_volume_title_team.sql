CREATE UNIQUE INDEX IF NOT EXISTS uniq_chapter_on_moderation_volume_title_team
ON chapters_on_moderation (lower(name), volume, title_id, team_id)
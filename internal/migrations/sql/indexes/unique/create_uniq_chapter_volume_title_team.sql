CREATE UNIQUE INDEX IF NOT EXISTS uniq_chapter_volume_title_team
ON chapters (lower(name), volume, title_id, team_id)
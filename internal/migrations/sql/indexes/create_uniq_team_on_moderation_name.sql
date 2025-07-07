CREATE UNIQUE INDEX IF NOT EXISTS uniq_team_on_moderation_name
ON teams_on_moderation (lower(name))
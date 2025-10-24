CREATE UNIQUE INDEX IF NOT EXISTS uniq_team_name
ON teams (lower(name))
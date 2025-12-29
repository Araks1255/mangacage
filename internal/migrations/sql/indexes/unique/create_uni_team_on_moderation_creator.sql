CREATE UNIQUE INDEX IF NOT EXISTS uni_teams_on_moderation_creator_id
ON teams_on_moderation (creator_id)
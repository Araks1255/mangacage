CREATE INDEX IF NOT EXISTS team_on_moderation_moderator_idx
ON teams_on_moderation (id, moderator_id)
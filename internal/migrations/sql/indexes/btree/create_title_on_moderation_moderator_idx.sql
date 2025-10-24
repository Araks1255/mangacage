CREATE INDEX IF NOT EXISTS title_on_moderation_moderator_idx
ON titles_on_moderation (id, moderator_id)
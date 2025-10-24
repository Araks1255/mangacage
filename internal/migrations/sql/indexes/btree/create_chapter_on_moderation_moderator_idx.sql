CREATE INDEX IF NOT EXISTS chapter_on_moderation_moderator_idx
ON chapters_on_moderation (id, moderator_id)
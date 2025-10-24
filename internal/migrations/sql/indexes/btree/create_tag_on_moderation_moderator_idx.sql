CREATE INDEX IF NOT EXISTS tag_on_moderation_moderator_idx
ON tags_on_moderation (id, moderator_id)
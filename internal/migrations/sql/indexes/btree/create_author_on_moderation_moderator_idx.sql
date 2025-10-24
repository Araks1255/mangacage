CREATE INDEX IF NOT EXISTS author_on_moderation_moderator_idx
ON authors_on_moderation (id, moderator_id)
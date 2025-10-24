CREATE INDEX IF NOT EXISTS genre_on_moderation_moderator_idx
ON genres_on_moderation (id, moderator_id)
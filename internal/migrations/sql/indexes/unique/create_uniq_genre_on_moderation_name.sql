CREATE UNIQUE INDEX IF NOT EXISTS uniq_genre_on_moderation_name
ON genres_on_moderation (lower(name))
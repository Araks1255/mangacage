CREATE UNIQUE INDEX IF NOT EXISTS uniq_genre_name
ON genres (lower(name))
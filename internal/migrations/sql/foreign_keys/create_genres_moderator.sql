DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_genres_moderator'
          AND table_name = 'genres'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE genres
        ADD CONSTRAINT fk_genres_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
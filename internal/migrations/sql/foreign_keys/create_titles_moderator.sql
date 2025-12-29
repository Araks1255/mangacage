DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_titles_moderator'
          AND table_name = 'titles'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE titles
        ADD CONSTRAINT fk_titles_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
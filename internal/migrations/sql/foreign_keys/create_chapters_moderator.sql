DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_chapters_moderator'
          AND table_name = 'chapters'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE chapters
        ADD CONSTRAINT fk_chapters_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_authors_moderator'
          AND table_name = 'authors'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE authors
        ADD CONSTRAINT fk_authors_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
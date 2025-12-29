DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_tags_moderator'
          AND table_name = 'tags'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE tags
        ADD CONSTRAINT fk_tags_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
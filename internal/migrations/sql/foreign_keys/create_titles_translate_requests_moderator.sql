DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_titles_translate_requests_moderator'
          AND table_name = 'titles_translate_requests'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE titles_translate_requests
        ADD CONSTRAINT fk_titles_translate_requests_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
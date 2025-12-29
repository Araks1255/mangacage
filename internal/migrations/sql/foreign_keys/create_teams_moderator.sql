DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_teams_moderator'
          AND table_name = 'teams'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE teams
        ADD CONSTRAINT fk_teams_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
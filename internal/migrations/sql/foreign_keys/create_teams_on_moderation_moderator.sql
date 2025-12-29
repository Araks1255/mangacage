DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_teams_on_moderation_moderator'
          AND table_name = 'teams_on_moderation'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE teams_on_moderation
        ADD CONSTRAINT fk_teams_on_moderation_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
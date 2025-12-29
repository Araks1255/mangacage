DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_users_moderator'
          AND table_name = 'users'
          AND constraint_type = 'FOREIGN KEY'
    ) THEN
        ALTER TABLE users
        ADD CONSTRAINT fk_users_moderator
        FOREIGN KEY (moderator_id)
        REFERENCES users (id) ON DELETE SET NULL;
    END IF;
END$$;
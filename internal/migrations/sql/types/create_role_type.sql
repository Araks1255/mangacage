DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type') THEN
        CREATE TYPE role_type AS ENUM (
            'team',
            'site'
        );
    ELSE
        ALTER TYPE role_type ADD VALUE IF NOT EXISTS 'team';
        ALTER TYPE role_type ADD VALUE IF NOT EXISTS 'site';
    END IF;
END$$;
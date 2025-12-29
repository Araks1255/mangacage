DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'title_translating_status') THEN
        CREATE TYPE title_translating_status AS ENUM (
            'ongoing',
            'completed',
            'abandoned',
            'suspended',
            'free'
        );
    ELSE
        ALTER TYPE title_translating_status ADD VALUE IF NOT EXISTS 'ongoing';
        ALTER TYPE title_translating_status ADD VALUE IF NOT EXISTS 'completed';
        ALTER TYPE title_translating_status ADD VALUE IF NOT EXISTS 'abandoned';
        ALTER TYPE title_translating_status ADD VALUE IF NOT EXISTS 'suspended';
        ALTER TYPE title_translating_status ADD VALUE IF NOT EXISTS 'free';
    END IF;
END$$;
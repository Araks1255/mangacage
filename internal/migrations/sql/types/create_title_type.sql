DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'title_type') THEN
        CREATE TYPE title_type AS ENUM (
            'manga',
            'manhwa',
            'manhua',
            'comics'
        );
    ELSE
        ALTER TYPE title_type ADD VALUE IF NOT EXISTS 'manga';
        ALTER TYPE title_type ADD VALUE IF NOT EXISTS 'manhwa';
        ALTER TYPE title_type ADD VALUE IF NOT EXISTS 'manhua';
        ALTER TYPE title_type ADD VALUE IF NOT EXISTS 'comics';
    END IF;
END$$;
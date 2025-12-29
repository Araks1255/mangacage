DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'title_publishing_status') THEN
        CREATE TYPE title_publishing_status AS ENUM (
            'announced',
            'ongoing',
            'completed',
            'suspended',
            'unknown'
        );
    ELSE
        ALTER TYPE title_publishing_status ADD VALUE IF NOT EXISTS 'announced';
        ALTER TYPE title_publishing_status ADD VALUE IF NOT EXISTS 'ongoing';
        ALTER TYPE title_publishing_status ADD VALUE IF NOT EXISTS 'completed';
        ALTER TYPE title_publishing_status ADD VALUE IF NOT EXISTS 'suspended';
        ALTER TYPE title_publishing_status ADD VALUE IF NOT EXISTS 'unknown';
    END IF;
END$$;
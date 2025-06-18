CREATE OR REPLACE FUNCTION set_title_translating_status_free()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM title_teams
        WHERE title_id = OLD.title_id
    ) THEN
         UPDATE titles
        SET translating_status = 'free'
        WHERE id = OLD.title_id;
    END IF;

    RETURN OLD;
END;
$$;

CREATE OR REPLACE FUNCTION increment_title_views()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles
    SET views = views + 1
    WHERE id = (
        SELECT v.title_id FROM volumes AS v
        INNER JOIN chapters AS c ON c.volume_id = v.id
        WHERE c.id = NEW.chapter_id
    );
    RETURN NEW;
END;
$$;
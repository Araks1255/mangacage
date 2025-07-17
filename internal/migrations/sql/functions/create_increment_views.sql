CREATE OR REPLACE FUNCTION increment_views()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE chapters
    SET views = views + 1
    WHERE id = NEW.chapter_id;

    UPDATE titles
    SET views = views + 1
    WHERE id = (
        SELECT title_id FROM chapters
        WHERE id = NEW.chapter_id
    );

    RETURN NEW;
END;
$$;
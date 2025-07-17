CREATE OR REPLACE FUNCTION increment_title_number_of_chapters()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles SET
    number_of_chapters = number_of_chapters + 1
    WHERE id = NEW.title_id;
    RETURN NEW;
END;
$$;
CREATE OR REPLACE FUNCTION decrement_title_number_of_chapters()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles SET
        number_of_chapters = number_of_chapters - 1
    WHERE
        id = OLD.title_id;
    RETURN OLD;
END;
$$;

CREATE OR REPLACE FUNCTION increment_chapter_on_moderation_number_of_pages()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE chapters_on_moderation SET
        number_of_pages = number_of_pages + 1
    WHERE
        id = NEW.chapter_on_moderation_id;
        
    RETURN NEW;
END;
$$;
CREATE OR REPLACE FUNCTION delete_chapters_on_moderation_pages_after_title_on_moderation_deletion() 
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM
        pages AS p
    USING
        chapters_on_moderation AS com
    WHERE  
        com.id = p.chapter_on_moderation_id
    AND
        com.title_on_moderation_id = OLD.id;
    
    RETURN OLD;
END;
$$;
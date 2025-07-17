CREATE OR REPLACE TRIGGER trg_decrement_title_number_of_chapters
AFTER DELETE ON chapters
FOR EACH ROW
EXECUTE FUNCTION decrement_title_number_of_chapters()
CREATE OR REPLACE TRIGGER trg_increment_title_number_of_chapters
AFTER INSERT ON chapters
FOR EACH ROW
EXECUTE FUNCTION increment_title_number_of_chapters()
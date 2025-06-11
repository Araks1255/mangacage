CREATE OR REPLACE TRIGGER trg_increment_title_views
AFTER INSERT ON user_viewed_chapters
FOR EACH ROW
EXECUTE FUNCTION increment_title_views()
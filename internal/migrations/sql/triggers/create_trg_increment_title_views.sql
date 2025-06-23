CREATE OR REPLACE TRIGGER trg_increment_views
AFTER INSERT ON user_viewed_chapters
FOR EACH ROW
EXECUTE FUNCTION increment_views()
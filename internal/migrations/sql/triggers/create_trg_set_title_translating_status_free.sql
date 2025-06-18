CREATE OR REPLACE TRIGGER trg_set_title_translating_status_free
AFTER DELETE ON title_teams
FOR EACH ROW
EXECUTE FUNCTION set_title_translating_status_free()
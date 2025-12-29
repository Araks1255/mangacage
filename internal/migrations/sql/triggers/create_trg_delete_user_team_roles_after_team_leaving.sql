CREATE OR REPLACE TRIGGER trg_delete_user_team_roles_after_team_leaving
AFTER UPDATE OF team_id ON users
FOR EACH ROW
EXECUTE FUNCTION delete_user_team_roles_after_team_leaving()
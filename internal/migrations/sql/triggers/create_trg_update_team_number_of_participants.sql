CREATE OR REPLACE TRIGGER trg_update_team_number_of_participants
AFTER UPDATE OF team_id ON users
FOR EACH ROW
EXECUTE FUNCTION update_team_number_of_participants()
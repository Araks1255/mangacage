CREATE OR REPLACE FUNCTION update_team_number_of_participants()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.team_id IS NOT NULL AND OLD.team_id IS NULL THEN
        UPDATE
            teams
        SET
            number_of_participants = number_of_participants + 1
        WHERE
            id = NEW.team_id;

    ELSIF NEW.team_id IS NULL AND OLD.team_id IS NOT NULL THEN
        UPDATE
            teams
        SET
            number_of_participants = number_of_participants - 1
        WHERE
            id = OLD.team_id;

    END IF;

    RETURN NEW;
END;
$$;

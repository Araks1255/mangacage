CREATE OR REPLACE FUNCTION delete_user_team_roles_after_team_leaving()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM
        user_roles AS ur
    USING
        roles AS r
    WHERE
        ur.role_id = r.id
    AND
        NEW.team_id IS NULL
    AND
        ur.user_id = NEW.id
    AND
        r.type = 'team'; 
    RETURN NEW;
END;
$$;
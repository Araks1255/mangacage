CREATE OR REPLACE FUNCTION fn_title_rates_after_update()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles SET
        sum_of_rates = sum_of_rates - OLD.rate + NEW.rate
    WHERE
        id = OLD.title_id;
    
    RETURN NEW;
END;
$$;
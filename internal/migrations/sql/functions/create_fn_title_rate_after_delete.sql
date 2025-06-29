CREATE OR REPLACE FUNCTION fn_title_rates_after_delete()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles SET
        sum_of_rates = sum_of_rates - OLD.rate,
        number_of_rates = number_of_rates - 1
    WHERE
        id = OLD.title_id;

    RETURN OLD;
END;
$$;
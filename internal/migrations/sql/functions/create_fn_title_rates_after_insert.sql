CREATE OR REPLACE FUNCTION fn_title_rates_after_insert()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE titles SET
        sum_of_rates = sum_of_rates + NEW.rate,
        number_of_rates = number_of_rates + 1
    WHERE
        id = NEW.title_id;

    RETURN NEW;
END;
$$;

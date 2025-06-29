CREATE OR REPLACE TRIGGER trg_title_rates_after_insert
AFTER INSERT ON title_rates
FOR EACH ROW
EXECUTE FUNCTION fn_title_rates_after_insert()